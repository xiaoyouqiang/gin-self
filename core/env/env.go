package env

import (
	"flag"
	"fmt"
	"strings"
)

type EnvType string

var (
	dev     = EnvType("dev")
	test    = EnvType("test")
	pre     = EnvType("pre")
	prod    = EnvType("prod")
	currEnv EnvType
)

//IsDev 是否是开发环境
func IsDev() bool {
	return currEnv == dev
}

//IsTest 是否是测试环境
func IsTest() bool {
	return currEnv == test
}

//IsPre 是否是预发布环境
func IsPre() bool {
	return currEnv == pre
}

//IsPro 是否是生成环境
func IsPro() bool {
	return currEnv == prod
}

//GetEnv 获得当前环境
func GetEnv() EnvType {
	return currEnv
}

func init() {
	env := flag.String("env", "", "清设置运行环境 dev | test | pre | prod\n")
	flag.Parse()

	e := strings.ToLower(strings.TrimSpace(*env))

	switch e {
	case "dev":
		currEnv = dev
	case "test":
		currEnv = test
	case "pre":
		currEnv = pre
	case "prod":
		currEnv = prod
	default:
		currEnv = dev
		fmt.Println("Warning: '-env' not found, or it is illegal. The default 'dev' will be used.")
	}
}