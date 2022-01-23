package middleware

/*
 设置 redis 上下文
*/

import (
	"gin-self/extend/self_redis"
	"github.com/gin-gonic/gin"
)

func SetRedisContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		self_redis.WithContext(c.Request.Context())
		c.Next()
	}
}