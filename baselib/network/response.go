package network

import (
	"net/http"
	"time"
	"encoding/json"
)

type Failurable interface {
	Code() string
	ErrorMsg() string
	DisplayedMsg() string
}

// 接口返回数据格式
type Response struct {
	HttpStatus  int         `json:"-"`
	Data        interface{} `json:"data"`
	Status      string      `json:"status"`
	ErrorCode   string      `json:"error_code"`
	ErrorMsg    string      `json:"error_msg"`
	DisErrorMsg string      `json:"dis_error_msg"`
	Timestamp   int64       `json:"timestamp"`
}

func NewFailure(f Failurable) *Response {
	return &Response{
		HttpStatus:  http.StatusOK,
		Status:      "fail",
		ErrorCode:   f.Code(),
		ErrorMsg:    f.ErrorMsg(),
		DisErrorMsg: f.DisplayedMsg(),
		Timestamp:   time.Now().Unix(),
	}
}

func NewSuccess(data interface{}) *Response {
	return &Response{
		HttpStatus: http.StatusOK,
		Status:     "success",
		Timestamp:  time.Now().Unix(),
		Data:       data,
	}
}

func NewHttpFailure(hs int, f Failurable) *Response {
	return &Response{
		HttpStatus:  hs,
		Status:      "fail",
		ErrorCode:   f.Code(),
		ErrorMsg:    f.ErrorMsg(),
		DisErrorMsg: f.DisplayedMsg(),
		Timestamp:   time.Now().Unix(),
	}
}

func (self *Response) AppendErr(errStr string) *Response {
	self.ErrorMsg += errStr
	return self
}

func (self *Response) Response(w http.ResponseWriter) {
	self.responseContext(w).response()
}

func (self *Response) responseContext(w http.ResponseWriter) *ResponseRecorder {
	if recorder, ok := w.(*ResponseRecorder); ok {
		if len(recorder.jsonBytes) != 0 {
			//todo
		}
		deBytes, _ := json.Marshal(self)
		recorder.jsonBytes = deBytes
		recorder.Status = self.HttpStatus
		return recorder
	} else {
		recorder = NewResponseRecorder(w)
		recorder.Status = self.HttpStatus
		deBytes, _ := json.Marshal(self)
		recorder.jsonBytes = deBytes
		return recorder
	}
}
