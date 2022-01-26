package middleware

/*
 全局获取用户信息
*/

import (
	"github.com/gin-gonic/gin"
)

func GetUserInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		//token := c.GetHeader("token")
		//
		//if token == "" {
		//	c.Abort()
		//}
		//
		//tokenData, err := my_jwt.ParseToken(token)
		//
		//if err != nil {
		//	c.Abort()
		//}
		//if _, ok := tokenData["data"]; !ok {
		//	c.Abort()
		//}
		//if _, ok := tokenData["data"].(map[string]interface{})["user_id"]; !ok {
		//	c.Abort()
		//}
		//if tokenData["data"].(map[string]interface{})["user_id"] == "" {
		//	c.Abort()
		//}

		//userIdString := tokenData["data"].(map[string]interface{})["user_id"]

		c.Set("user", gin.H{})

		c.Next()
	}
}

//ValueUser 获取用户信息
func ValueUser(c *gin.Context) map[string]interface{} {
	if user, ok := c.Get("user"); ok {
		if user != nil {
			return user.(map[string]interface{})
		}
	}

	return map[string]interface{}{}
}
