package helpers

import (
	"gin-self/extend/e"
	"gin-self/extend/self_loger"

	"github.com/gin-gonic/gin"
)

func ApiSuccess(c *gin.Context, data interface{}) {
	code := e.SUCCESS
	traceId := self_loger.GetTraceByCtx(c).ValueTraceId()
	c.JSON(code, gin.H{"code": code, "message": e.GetMessage(code), "data": data, "trace_id": traceId})
}

//ApiError 错误返回
//message 错误消息可以重写 非必传
//demo 	helpers.ApiError(c, e.AUTH_FAIL, "new message")
func ApiError(c *gin.Context, code int, message ...string) {
	var newMessage string
	for _, v := range message {
		newMessage = v
		break
	}
	if newMessage == "" {
		newMessage = e.GetMessage(code)
	}

	traceId := self_loger.GetTraceByCtx(c).ValueTraceId()
	c.JSON(e.SUCCESS, gin.H{"code": code, "message": newMessage, "data": gin.H{}, "trace_id": traceId})
}
