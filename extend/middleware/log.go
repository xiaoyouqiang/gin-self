package middleware

import (
	"bytes"
	"context"
	"time"

	"gin-self/extend/self_loger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//拦截gin 响应数据 记录到日志
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LogMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := self_loger.GetInstance()
		trace := self_loger.NewTrace()

		childCtx := context.WithValue(c.Request.Context(),"trace", trace)

		c.Request = c.Request.WithContext(childCtx)

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}

		c.Writer = blw

		startTime := time.Now()

		_ = c.Request.ParseForm()

		c.Next()

		//执行耗时
		endTime := time.Now()
		execTime := endTime.Sub(startTime).Seconds()

		//body, _ := c.GetRawData()

		traceInfo := logrus.Fields {
			"trace_id":    trace.ValueTraceId(),
			"status_code": c.Writer.Status(),
			"host":        c.Request.Host,
			"req_method":  c.Request.Method,
			"req_uri":     c.Request.RequestURI,
			"exec_time":   execTime,
			"client_ip":   c.ClientIP(),
			"trace_info": logrus.Fields {
				//"header_params": c.Request.Header, //c.GetHeader("Cookie"),后续着重记录一些特定的头，不用记录所有header
				"post_params"	:	c.Request.PostForm,
				"response":      blw.body.String(),
				"sql":           trace.ValueSqlInfo(),
				"redis":         trace.ValueRedisInfo(),
				"error_stack":    trace.ValueErrorInfo(),
			},
		}

		self_loger.ReleaseTrace(trace)

		logger.Info(traceInfo)
	}
}
