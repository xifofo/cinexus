package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 标准API响应结构
type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

// Success 返回成功响应
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "成功",
		Data: data,
	})
}

// SuccessWithMsg 返回带自定义消息的成功响应
func SuccessWithMsg(c *gin.Context, msg string, data any) {
	c.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  msg,
		Data: data,
	})
}

// Fail 返回失败响应
func Fail(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
	})
}

// BadRequest 返回400错误响应
func BadRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, Response{
		Code: 400,
		Msg:  msg,
	})
}

// Unauthorized 返回401错误响应
func Unauthorized(c *gin.Context, msg string) {
	if msg == "" {
		msg = "未授权"
	}
	c.JSON(http.StatusUnauthorized, Response{
		Code: 401,
		Msg:  msg,
	})
}

// Forbidden 返回403错误响应
func Forbidden(c *gin.Context, msg string) {
	if msg == "" {
		msg = "禁止访问"
	}
	c.JSON(http.StatusForbidden, Response{
		Code: 403,
		Msg:  msg,
	})
}

// NotFound 返回404错误响应
func NotFound(c *gin.Context, msg string) {
	if msg == "" {
		msg = "资源不存在"
	}
	c.JSON(http.StatusNotFound, Response{
		Code: 404,
		Msg:  msg,
	})
}

// ServerError 返回500错误响应
func ServerError(c *gin.Context, msg string) {
	if msg == "" {
		msg = "服务器内部错误"
	}
	c.JSON(http.StatusInternalServerError, Response{
		Code: 500,
		Msg:  msg,
	})
}
