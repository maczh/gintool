package gintool

import (
	"gopkg.in/mgo.v2/bson"
)

type PostLog struct {
	ID           bson.ObjectId          `bson:"_id"`
	Time         string                 `json:"time" bson:"time"`
	RequestId    string                 `json:"requestId" bson:"requestId"`
	Responsetime string                 `json:"responsetime" bson:"responsetime"`
	TTL          int                    `json:"ttl" bson:"ttl"`
	Apiname      string                 `json:"apiName" bson:"apiName"`
	Controller   string                 `json:"controller" bson:"controller"`
	Tradenum     string                 `json:"tradenum" bson:"tradenum"`
	Accountid    string                 `json:"accountid" bson:"accountid"`
	Requestparam map[string]string      `json:"requestparam" bson:"requestparam"`
	Responsestr  string                 `json:"responsestr" bson:"responsestr"`
	Responsemap  map[string]interface{} `json:"responsemap" bson:"responsemap"`
}
