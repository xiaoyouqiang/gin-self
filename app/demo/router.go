package demo

import (
	"github.com/gin-gonic/gin"
)

func ReturnRoute(r *gin.Engine) {
	//该方法需要登录验证， token检查中间件: middleware.CheckToken
	r.Any("/demo/test/index", Index)
	//r.GET("/demo/test/index1",middleware.CheckTokenMiddleWare(), Index1)
	r.GET("/demo/test/index1", Index1)
	r.GET("/demo/test/index2", Index2)
	r.POST("/demo/test/index2", Index2)
	r.GET("/demo/test/redis_test", RedisTest)
	r.GET("/demo/test/http_test", HttpTest)
}
