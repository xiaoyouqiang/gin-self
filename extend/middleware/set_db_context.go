package middleware

/*
 设置 db 上下文
*/

import (
	"gin-self/extend/self_db"
	"github.com/gin-gonic/gin"
)

func SetDbContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		self_db.WithContext(c.Request.Context())
		c.Next()
	}
}