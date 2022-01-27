package helpers

import (
	"gin-self/extend/self_loger"
	"gin-self/extend/utils/e"

	"github.com/gin-gonic/gin"
)

func ApiSuccess(c *gin.Context, data interface{}) {
	traceId := self_loger.GetTraceByCtx(c).ValueTraceId()
	c.JSON(e.Success, gin.H{"code": e.Success, "message": e.GetMessage(e.Success), "data": data, "trace_id": traceId})
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
	c.JSON(e.Success, gin.H{"code": code, "message": newMessage, "data": gin.H{}, "trace_id": traceId})
}
