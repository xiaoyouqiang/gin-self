package middleware

/*
 token检查中间件
*/

import (
	"gin-self/extend/e"
	"gin-self/extend/helpers"
	"gin-self/extend/my_jwt"
	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
)

//不需要检查登录的url
var notNeedCheckLoginUrl = map[string]struct{}{
	"/demo/test/http_test": {},
	"/demo/test/index": {},
}

//检查是否需要登录
func isNeedLogin(c *gin.Context) bool {
	url := c.FullPath()
	if _,ok := notNeedCheckLoginUrl[url];ok {
		//不需要登录
		return false
	}
	//需要登录
	return true
}

func CheckToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")

		if token == "" && isNeedLogin(c) {
			helpers.ApiError(c, e.AUTH_FAIL)
			c.Abort()
			return
		}

		var tokenData = jwt.MapClaims{}
		var err error
		if token != "" {
			tokenData, err = my_jwt.ParseToken(token)
		}

		if err != nil && isNeedLogin(c) {
			helpers.ApiError(c, e.AUTH_FAIL)
			c.Abort()
			return
		}
		if _, ok := tokenData["data"]; !ok && isNeedLogin(c)  {
			helpers.ApiError(c, e.AUTH_FAIL)
			c.Abort()
			return
		}
		if isNeedLogin(c) && tokenData["data"].(map[string]interface{})["user_id"] == "" {
			helpers.ApiError(c, e.AUTH_FAIL)
			c.Abort()
			return
		}
		if userId, ok := tokenData["data"]; ok && tokenData["data"].(map[string]interface{})["user_id"] != "" {
			//获取用户信息
			SetLoginUserInfo(c, userId.(int))
		}
	}
}

type LoginUserInfo struct {
	Id int
	UserName string
}

//GetLoginUserInfo 获取登录用户信息
func GetLoginUserInfo(ctx *gin.Context) LoginUserInfo {
	if v,ok := ctx.Get("user_info");ok {
		return v.(LoginUserInfo)
	}

	return LoginUserInfo{}
}

//SetLoginUserInfo 设置登录用户信息
func SetLoginUserInfo(ctx *gin.Context, userId int) {
	//查询数据库

	//查询到的信息保持到 ctx 中
	ctx.Set("user_info", LoginUserInfo{})
}

