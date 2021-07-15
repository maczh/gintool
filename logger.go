package gintool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/maczh/logs"
	"github.com/maczh/mgconfig"
	"github.com/maczh/mgtrace"
	"github.com/maczh/utils"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

var accessChannel = make(chan string, 100)

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func SetRequestLogger() gin.HandlerFunc {

	go handleAccessChannel()

	return func(c *gin.Context) {
		bodyLogWriter := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bodyLogWriter

		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		responseBody := bodyLogWriter.body.String()

		var result map[string]interface{}

		// 日志格式
		if strings.Contains(c.Request.RequestURI, "/docs") {
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
		params := utils.GinParamMap(c)
		postLog := new(PostLog)
		postLog.ID = bson.NewObjectId()
		postLog.Time = startTime.Format("2006-01-02 15:04:05")
		if strings.Contains(c.Request.RequestURI, "?") {
			postLog.Controller = c.Request.RequestURI[0:strings.Index(c.Request.RequestURI, "?")]
		} else {
			postLog.Controller = c.Request.RequestURI
		}
		postLog.RequestId = mgtrace.GetRequestId()
		postLog.Requestparam = params
		postLog.Responsetime = endTime.Format("2006-01-02 15:04:05")
		postLog.Responsemap = result
		postLog.TTL = int(endTime.UnixNano()/1e6 - startTime.UnixNano()/1e6)

		accessLog := "|" + c.Request.Method + "|" + postLog.Controller + "|" + c.ClientIP() + "|" + endTime.Format("2006-01-02 15:04:05.012") + "|" + fmt.Sprintf("%vms", endTime.UnixNano()/1e6-startTime.UnixNano()/1e6)
		logs.Debug(accessLog)
		logs.Debug("请求参数:{}", postLog.RequestId, params)
		logs.Debug("接口返回:{}", result)

		if mgconfig.GetConfigString("go.log.req") != "" {
			accessChannel <- utils.ToJSON(postLog)
		}
	}
}

func handleAccessChannel() {
	for accessLog := range accessChannel {
		var postLog PostLog
		json.Unmarshal([]byte(accessLog), &postLog)
		err := mgconfig.Mgo.C(mgconfig.GetConfigString("go.log.req")).Insert(postLog)
		if err != nil {
			logs.Error("MongoDB写入错误:" + err.Error())
		}
	}
	return
}
