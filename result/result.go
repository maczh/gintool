package result

import "strconv"

/**
通用返回结果类
*/
type Result struct {
	Status int32       `json:"status" bson:"status"`
	Msg    string      `json:"msg" bson:"msg"`
	Data   interface{} `json:"data" bson:"data"`
	Page   *ResultPage `json:"page" bson:"page"`
}

type ResultPage struct {
	Count int `json:"count"` //总页数
	Index int `json:"index"` //页号
	Size  int `json:"size"`  //分页大小
	Total int `json:"total"` //总记录数
}

func Success(d interface{}) *Result {
	result := new(Result)
	result.Data = d
	result.Status = 1
	result.Msg = "成功"
	return result
}

func SuccessWithPage(d interface{}, count, index, size, total int) *Result {
	result := new(Result)
	result.Data = d
	result.Status = 1
	result.Msg = "成功"
	page := new(ResultPage)
	page.Count = count
	page.Index = index
	page.Size = size
	page.Total = total
	result.Page = page
	return result
}

func Error(s int32, m string) *Result {
	result := new(Result)
	result.Status = s
	result.Msg = m
	return result
}

type AppResult struct {
	Status string      `json:"status" bson:"status"`
	Msg    string      `json:"msg" bson:"msg"`
	Data   interface{} `json:"data" bson:"data"`
	Page   *ResultPage `json:"page" bson:"page"`
}

func NewAppResult(r Result) AppResult {
	return AppResult{
		Status: strconv.Itoa(int(r.Status)),
		Msg:    r.Msg,
		Data:   r.Data,
		Page:   r.Page,
	}
}

func AppError(s int, m string) AppResult {
	return AppResult{
		Status: strconv.Itoa(s),
		Msg:    m,
		Data:   nil,
		Page:   nil,
	}
}
