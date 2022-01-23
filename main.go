package main

import (
	"gin-self/core/app"
	"gin-self/extend/self_loger"
)

func main() {
	application := app.New()

	application.InitMysql()

	application.InitRedis()

	application.InitLogger(self_loger.HttpType)

	application.Run()
}
