package http

import (
	"github.com/gin-gonic/gin"
)

const (
	// OK ok
	OK = 200
	// RequestErr request error
	RequestErr = 400
	ForbiddenErr = 403
	// ServerErr server error
	ServerErr    = 500


	contextErrCode      = "context/err/code"
	ParamIncorrect      = "ParamIncorrect"
	ConfigErr           = "Config error"
	MemberAuthIncorrect = "MemberAuthIncorrect"
)

type resp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func errors(c *gin.Context, code int, msg string) {
	c.Set(contextErrCode, code)
	c.JSON(200, resp{
		Code:    code,
		Message: msg,
	})

	return
}

func result(c *gin.Context, data interface{}, code int) {
	c.Set(contextErrCode, code)
	c.JSON(200, resp{
		Code: code,
		Data: data,
	})
}

func handleResult(c *gin.Context, err error, data interface{}) {
	if err != nil {
		errors(c, RequestErr, err.Error())
		return
	}
	result(c, data, OK)
	return
}

func handleResultWithCode(c *gin.Context, err error, code int, data interface{}) {
	if err != nil {
		errors(c, code, err.Error())
		return
	}
	result(c, data, OK)
	return
}

func failResult(c *gin.Context, code int, msg string) {
	c.Set(contextErrCode, code)
	c.JSON(200, resp{
		Code:    code,
		Message: msg,
	})

	return
}

func successResult(c *gin.Context, data interface{}) {
	c.Set(contextErrCode, 200)
	c.JSON(200, resp{
		Code: OK,
		Data: data,
	})
}
