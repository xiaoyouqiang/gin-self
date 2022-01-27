package app

import (
	"gin-self/app/demo"
	"gin-self/core/env"
	"gin-self/extend/middleware"
	"gin-self/extend/utils/e"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func IncludeRoute() *gin.Engine {
	router := gin.New()
	router.Use(middleware.LogMiddleWare())
	router.Use(middleware.TimeOut())
	router.Use(middleware.SetDbContext())
	router.Use(middleware.SetRedisContext())
	//router.Use(middleware.CheckToken())
	//router.Use(middleware.CheckSign())
	router.Use(middleware.Recovery())

	if !env.IsPro() {
		pprof.Register(router) // register pprof to gin
	}

	gin.SetMode(gin.ReleaseMode)

	//404 handler
	router.NoRoute(func(c *gin.Context) {
		code := e.NotFound
		c.JSON(404, gin.H{"code": code, "message": e.GetMessage(code)})
	})

	//加载app/demo模块的独立路由，不同模块各自维护自己的路由配置
	demo.ReturnRoute(router)

	return router
}
