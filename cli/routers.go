package cli

import (
	"gin-self/cli/test"
	"gin-self/core"
)

func IncludeRoute(r *core.CliEngine) {
	test.ReturnRoute(r)
}
