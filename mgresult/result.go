package mgresult

import "strconv"

/**
通用返回结果类
*/
type Result struct {
	Status int         `json:"status" bson:"status"`
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

func Success(d interface{}) Result {
	return Result{
		Status: 1,
		Msg:    "成功",
		Data:   d,
	}
}

func SuccessWithPage(d interface{}, count, index, size, total int) Result {
	result := Result{
		Status: 1,
		Msg:    "成功",
		Data:   d,
		Page: &ResultPage{
			Count: count,
			Index: index,
			Size:  size,
			Total: total,
		},
	}
	return result
}

func Error(s int, m string) Result {
	return Result{
		Status: s,
		Msg:    m,
	}
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
