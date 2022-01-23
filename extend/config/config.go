package config

import (
	"gin-self/core/env"
	"log"

	"github.com/go-ini/ini"
)

var config *ini.File

func init() {
	var err error
	source := "conf/" + string(env.GetEnv()) + "/app.ini"
	config, err = ini.LoadSources(ini.LoadOptions{UnescapeValueDoubleQuotes: true}, source)
	if err != nil {
		log.Fatalf("Fail to parse ini': %v", err)
	}
}

// Get
// @description   读取ini配置
// @auth      xiao you qiang
// @param     section        string         要读取的 ini 文章的章节 没有配置章节的传 空字符串
// @param     key        string         ini 章节下 要读取的配置key
// @return    Key        ini.key        返回 ini包key的引用 调用方可以继续 key.*一些列操作
func Get(section string, key string) *ini.Key {
	return config.Section(section).Key(key)
}

//Section 配置章节是否存在
func Section(name string) *ini.Section {
	return config.Section(name)
}
