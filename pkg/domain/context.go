package domain

import (
	"net/http"
)

type Context interface {
	Request() *http.Request
	Response() http.ResponseWriter
}

type context struct {
	request  *http.Request
	response http.ResponseWriter
}

func NewContext(request *http.Request, response http.ResponseWriter) Context {
	return &context{
		request:  request,
		response: response,
	}
}

func (c *context) Request() *http.Request {
	return c.request
}

func (c *context) Response() http.ResponseWriter {
	return c.response
}
