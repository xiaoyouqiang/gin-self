package app

import (
	"flag"
	"gin-self/cli"
	"gin-self/core"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
	"log"
	"time"
)

type CliApp struct {
	ctx    context.Context
	cancel func()
	App
}

func NewCli() *CliApp {
	ctx, cancel := context.WithCancel(context.Background())
	return &CliApp{
		ctx: ctx,
		cancel: cancel,
	}
}

func (a *CliApp) RunCli() {
	group,_ := errgroup.WithContext(a.ctx)
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		log.Println("cli format error,like this")
		log.Println("cli command_name params, params must not required")
		log.Fatalln("when if have params like this a=1 b=2")
	}

	engine := core.NewCLiEngine(args)

	cli.IncludeRoute(engine)

	group.Go(func() (err error) {
		engine.ExecHandler()
		return
	})

	log.Printf("%s: cli exit %v\n", time.Now().Format("2006-01-02 15:04:05"), group.Wait())
}