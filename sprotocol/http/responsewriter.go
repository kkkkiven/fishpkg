package http

import (
	"fmt"
	"net/http"
)

//go:generate msgp

type ResponseWriter struct {
	Body          []byte              `msg:"body"`
	HandlerHeader map[string][]string `msg:"header"`
	Status        int                 `msg:"code"`
}

func NewResponseWriter() *ResponseWriter {
	rw := new(ResponseWriter)
	rw.Body = make([]byte, 0, 1024)
	rw.HandlerHeader = make(map[string][]string, 0)
	return rw
}

func (this *ResponseWriter) Write(p []byte) (int, error) {
	this.Body = append(this.Body, p...)
	return len(p), nil
}

func (this *ResponseWriter) Header() http.Header {
	return this.HandlerHeader
}

func (this *ResponseWriter) WriteHeader(statusCode int) {
	this.checkWriteHeaderCode(statusCode)

	this.Status = statusCode
}

func (this *ResponseWriter) checkWriteHeaderCode(code int) {
	if code != 0 {
		if code < 100 || code > 999 {
			panic(fmt.Sprintf("invalid header code: %v", code))
		}
	}
}
