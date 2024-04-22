package domain

import (
	"net/url"

	"github.com/gin-gonic/gin"
)

type H = map[string]interface{}

// ResponseStruct Http response builder
type ResponseStruct struct {
	status int
	body   H
}

func (h *ResponseStruct) Data(data H) *ResponseStruct {
	h.body["data"] = data
	return h
}

func (h *ResponseStruct) Errors(e H) *ResponseStruct {
	h.body["errors"] = e
	return h
}

func (h *ResponseStruct) Message(m string) *ResponseStruct {
	h.body["message"] = m
	return h
}

func (h *ResponseStruct) Send(c *gin.Context) {
	c.JSON(h.status, h.body)
}

func (h *ResponseStruct) Redirect(c *gin.Context, path string) {
	absolutePath := "/"
	u := url.URL{Path: absolutePath}
	c.Request.URL = &u
	c.Redirect(h.status, path)
}

func HttpResponse(statusCode int) *ResponseStruct {
	return &ResponseStruct{
		status: statusCode,
		body: H{
			"status_code": statusCode,
		},
	}
}
