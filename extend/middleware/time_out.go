package middleware

/*
 全局超时控制
*/

import (
	"context"
	"gin-self/extend/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func TimeOut() gin.HandlerFunc {
	return func(c *gin.Context) {
		timeOut := config.Get("server", "global_timeout").MustInt()
		ctx, cancel := context.WithTimeout(c.Request.Context(), time.Duration(timeOut) * time.Second)
		defer func() {
			cancel()
			if ctx.Err() == context.DeadlineExceeded {
				c.Writer.WriteHeader(http.StatusGatewayTimeout)
				c.Abort()
			}
		}()

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}