package gintool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maczh/logs"
	"github.com/maczh/mgconfig"
	"github.com/maczh/mgtrace"
	"github.com/maczh/utils"
	"gopkg.in/mgo.v2/bson"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

var accessChannel = make(chan string, 100)
var collection string

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func SetRequestLogger() gin.HandlerFunc {

	if collection == "" {
		collection = mgconfig.GetConfigString("go.log.req")
	}
	go handleAccessChannel()

	return func(c *gin.Context) {
		bodyLogWriter := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bodyLogWriter

		// 开始时间
		startTime := time.Now()

		data, err := c.GetRawData()
		if err != nil {
			logs.Error("GetRawData error:", err.Error())
		}
		body := string(data)

		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data)) // 关键点

		// 处理请求
		c.Next()

		responseBody := bodyLogWriter.body.String()

		var req map[string]interface{}
		var result map[string]interface{}

		// 日志格式
		if strings.Contains(c.Request.RequestURI, "/docs") || c.Request.RequestURI == "/" {
			return
		}

		if responseBody != "" && responseBody[0:1] == "{" {
			err := json.Unmarshal([]byte(responseBody), &result)
			if err != nil {
				result = map[string]interface{}{"status": -1, "msg": "解析异常:" + err.Error()}
			}
		}

		// 结束时间
		endTime := time.Now()

		// 日志格式
		var params interface{}
		if strings.Contains(c.ContentType(), "application/json") {
			utils.FromJSON(body, &req)
			params = req
		} else if strings.Contains(c.ContentType(), "x-www-form-urlencoded") || strings.Contains(c.ContentType(), "multipart/form-data") {
			params = utils.GinParamMap(c)
		} else if c.Request.Method != "GET" && c.Request.Method != "DELETE" {
			return
		}
		postLog := new(PostLog)
		postLog.ID = bson.NewObjectId()
		postLog.Time = startTime.Format("2006-01-02 15:04:05")
		postLog.Uri = c.Request.RequestURI
		postLog.Method = c.Request.Method
		postLog.RequestId = mgtrace.GetRequestId()
		postLog.ContentType = c.ContentType()
		postLog.RequestHeader = getHeaders(c)
		ip := c.GetHeader("X-Forward-For")
		if ip == "" {
			ip = c.GetHeader("X-Real-IP")
			if ip == "" {
				ip = c.ClientIP()
			}
		}
		postLog.ClientIP = ip
		postLog.RequestParam = params
		postLog.ResponseTime = endTime.Format("2006-01-02 15:04:05")
		postLog.ResponseMap = result
		postLog.TTL = int(endTime.UnixNano()/1e6 - startTime.UnixNano()/1e6)

		accessLog := "|" + c.Request.Method + "|" + postLog.Uri + "|" + c.ClientIP() + "|" + endTime.Format("2006-01-02 15:04:05.012") + "|" + fmt.Sprintf("%vms", endTime.UnixNano()/1e6-startTime.UnixNano()/1e6)
		logs.Debug(accessLog)
		logs.Debug("请求参数:{}", params)
		logs.Debug("接口返回:{}", result)

		if collection != "" || mgconfig.GetConfigString("go.log.kafka.use") == "true" {
			accessChannel <- utils.ToJSON(postLog)
		}
	}
}

func handleAccessChannel() {
	for accessLog := range accessChannel {
		if mgconfig.GetConfigString("go.log.kafka.use") == "true" && mgconfig.GetConfigString("go.log.kafka.topic") != "" && strings.Contains(mgconfig.GetConfigString("go.config.used"), "kafka") {
			topics := strings.Split(mgconfig.GetConfigString("go.log.kafka.topic"), ",")
			for _, topic := range topics {
				err := mgconfig.Kafka.Send(topic, accessLog)
				if err != nil {
					logs.Error("接口日志发送到kafka的{}主题失败:{}", topic, err.Error())
				}
			}
		}
		if collection == "" {
			continue
		}
		var postLog PostLog
		json.Unmarshal([]byte(accessLog), &postLog)
		mgo, err := mgconfig.GetMongoConnection()
		if err != nil {
			logs.Error("MongoDB连接错误:{}", err.Error())
			continue
		}
		err = mgo.C(collection).Insert(postLog)
		if err != nil {
			logs.Error("MongoDB写入错误:" + err.Error())
		}
		mgconfig.ReturnMongoConnection(mgo)
	}
	return
}

func getHeaders(c *gin.Context) map[string]string {
	headers := make(map[string]string)
	for k, v := range c.Request.Header {
		headers[k] = v[0]
	}
	return headers
}
