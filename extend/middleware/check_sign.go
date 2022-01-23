package middleware

/*
 签名校验中间件
*/

import (
	"crypto/md5"
	"fmt"
	"gin-self/conf"
	"gin-self/extend/e"
	"gin-self/extend/helpers"
	"github.com/gin-gonic/gin"
	"sort"
	"strings"
)

var (
	//需要参与签名的头部
	needSignHeader = map[string]struct{} {
		"appid":{},
		"appkey":{},
		"timestamp":{},
		"apiplateform":{},
		"apiappbundleid":{},
		"noncestr":{},
	}
)

func CheckSign() gin.HandlerFunc {
	return func(c *gin.Context) {
		//检查签名头部都提交
		var headerData = map[string]string{}
		for key,_ := range needSignHeader {
			if c.GetHeader(key) == "" {
				helpers.ApiError(c, e.SIGN_ERROR,"not enough sign param")
				c.Abort()
				return
			}
			headerData[key] = c.GetHeader(key)
		}

		//获取 appSecret
		appSecret,ok := getAppSecret(c.GetHeader("appid"), c.GetHeader("appkey"))
		if !ok {
			helpers.ApiError(c, e.SIGN_ERROR,"app key error")
			c.Abort()
			return
		}

		//获取get或post参数 保留需要参与签名的数据
		signData := getNeedSignParams(c)

		//合并签名数据
		newSignData := mergeDataForMap(headerData, signData)

		//签名数组排序
		needSignStr := makeSignStrForSort(newSignData)

		//根据不同端进行签名计算
		serverSign := makeSign(needSignStr,appSecret, c.GetHeader("appkey"))

		//校验提交的签名与服务器生成的签名
		if serverSign != c.GetHeader("sign") {
			helpers.ApiError(c, e.SIGN_ERROR,"sign error")
			c.Abort()
			return
		}

		c.Next()
	}
}

func getAppSecret(appId, appkey string) (string,bool)  {
	if v,ok := conf.AppKeySecret[appId][appkey];ok {
		return v,true
	}

	return "",false
}

func getNeedSignParams(ctx *gin.Context) map[string]string {
	var data map[string][]string
	var needSignData = map[string]string{}

	if ctx.Request.Method == "GET" {
		data = ctx.Request.URL.Query()
	}
	if ctx.Request.Method == "POST" {
		data = ctx.Request.PostForm
	}

	for k,v := range data {
		if len(v) > 1 {
			//数组不参与签名
			continue
		}

		if v[0] == "" || v[0] == "0" || strings.ToLower(v[0]) == "null" {
			//空串不参与签名
			continue
		}
		if k == "sign" {
			//sign不参与签名
			continue
		}

		needSignData[k] = v[0]
	}

	return needSignData
}

func mergeDataForMap(mapData ...map[string]string) map[string]string {
	var newMap = map[string]string{}
	for _,m := range mapData {
		for k,v := range m {
			newMap[k] = v
		}
	}

	return newMap
}

func makeSignStrForSort(data map[string]string) string {
	var keys []string
	var signStrArray []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _,k := range keys {
		signStrArray = append(signStrArray, k + "=" + data[k])
	}

	return strings.Join(signStrArray,",")
}

func makeSign(signsTr,appSecret, appKey string) string  {
	//app端签名计算
	var s string
	if strings.Contains(appKey, "IOS") || strings.Contains(appKey, "ANDROID") {
		s = fmt.Sprintf("%x", md5.Sum([]byte(appSecret + signsTr + appSecret)))
	} else {
		//非app端签名计算

		b := fmt.Sprintf("%x",md5.Sum([]byte(appSecret + signsTr + appSecret)))

		s = fmt.Sprintf("%x", md5.Sum([]byte(b)))
	}

	return s
}
