package main

import (
	"fmt"
	"net/http"
)

type ResponseWriter struct {
	resp   http.ResponseWriter
	status int
	size   int64
}

func NewResponseWriter(resp http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		resp,
		http.StatusOK,
		0,
	}
}

func (r *ResponseWriter) Header() http.Header {
	return r.resp.Header()
}

func (r *ResponseWriter) Write(b []byte) (int, error) {
	r.size += int64(len(b))
	return r.resp.Write(b)
}

func (r *ResponseWriter) WriteHeader(statusCode int) {
	r.status = statusCode
	r.resp.WriteHeader(statusCode)
}

func (r *ResponseWriter) Status() int {
	return r.status
}

func (r *ResponseWriter) Size() int64 {
	return r.size
}

func (r *ResponseWriter) WriteError(msg string) {
	r.resp.Header().Set("Content-Type", "application/json")
	r.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(r.resp, "{\"__type__\":\"Error\",\"message\":\""+msg+"\"}")
}
