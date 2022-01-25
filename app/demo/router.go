package demo

import (
	"github.com/gin-gonic/gin"
)

func ReturnRoute(r *gin.Engine) {
	r.Any("/demo/test/index", Index)
	r.GET("/demo/test/index1", Index1)
	r.GET("/demo/test/index2", Index2)
	r.POST("/demo/test/index2", Index2)
	r.GET("/demo/test/redis_test", RedisTest)
	r.GET("/demo/test/http_test", HttpTest)
}
