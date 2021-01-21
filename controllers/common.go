package controllers

import (
	"math/rand"
	"time"

	"github.com/astaxie/beego"
)

type CommonController struct {
	beego.Controller
}

type JsonStruct struct {
	Code  int         `json:"code"`
	Msg   interface{} `json:"msg"`
	Items interface{} `json:"items"`
	Count int64       `json:"count"`
}

func ReturnSuccess(code int, msg interface{}, items interface{}, count int64) (json *JsonStruct) {
	json = &JsonStruct{Code: code, Msg: msg, Items: items, Count: count}
	return
}

func ReturnError(code int, msg interface{}) (json *JsonStruct) {
	json = &JsonStruct{Code: code, Msg: msg}
	return
}


func DateFormat(times int64) string {
	videoTime := time.Unix(times, 0)
	return videoTime.Format("2006-01-02")
}


func RandString(length int) string {
	var source = rand.NewSource(time.Now().UnixNano())

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[source.Int63()%int64(len(charset))]
	}
	return string(b)
}