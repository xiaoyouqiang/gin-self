package test

import (
	"gin-self/core"
)

func ReturnRoute(r *core.CliEngine) {
	//r.GET("/demo/test/index1",middleware.CheckTokenMiddleWare(), Index1)
	r.SetHandler("test.test", Test)
}
