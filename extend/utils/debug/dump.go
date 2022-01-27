package debug

import (
	"gin-self/extend/utils/e"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
)

//VarDump 控制台输出
func VarDump(a ...interface{}) {
	spew.Dump(a...)
}

//HttpVarDump 输出到 http
func HttpVarDump(c *gin.Context, a ...interface{}) {
	c.String(e.Success, spew.Sdump(a...))
}