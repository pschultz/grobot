package testAPI

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func NewHttpResponse(body string) *http.Response {
	responseBody := ioutil.NopCloser(bytes.NewBuffer([]byte(body)))
	contentLength := int64(len(body))
	return &http.Response{
		Status:        "200 OK",
		StatusCode:    200,
		Proto:         "HTTP/1.0",
		ProtoMajor:    1,
		ProtoMinor:    0,
		Header:        make(map[string][]string),
		Body:          responseBody,
		ContentLength: contentLength,
	}
}
