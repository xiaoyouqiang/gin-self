package debug

import (
	"github.com/davecgh/go-spew/spew"
)

func VarDump(a ...interface{}) {
	spew.Dump(a...)
}
