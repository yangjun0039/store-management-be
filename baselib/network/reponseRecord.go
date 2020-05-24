package network

import "net/http"

type ResponseRecorder struct {
	Status    int
	uuid      string
	jsonBytes []byte
	w         http.ResponseWriter
}

func NewResponseRecorder(w http.ResponseWriter) *ResponseRecorder {
	return &ResponseRecorder{Status: 200, jsonBytes: []byte{}, w: w}
}

func (self *ResponseRecorder) Header() http.Header {
	return self.w.Header()
}

func (self *ResponseRecorder) Write(b []byte) (int, error) {
	return self.w.Write(b)
}

func (self *ResponseRecorder) WriteHeader(status int) {
	self.w.WriteHeader(status)
}

func (self *ResponseRecorder) response() {
	self.WriteHeader(self.Status)
	self.Write(self.jsonBytes)
	// todo 日志记录
}

