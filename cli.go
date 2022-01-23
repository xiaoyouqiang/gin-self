package main

import (
	"gin-self/core/app"
	"gin-self/extend/self_loger"
)

func main() {
	application := app.NewCli()

	application.InitMysql()

	application.InitRedis()

	application.InitLogger(self_loger.CliType)

	application.RunCli()
}
