package self_jwt

import (
	"fmt"
	"time"

	"gin-self/extend/config"

	jwt "github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte(config.Get("jwt", "secret").MustString(""))

type TokenStruct struct {
	userId     int64 `json:"user_id"`
	updateTime int64 `json:"update_time"`
	jwt.StandardClaims
}

//CreateToken 生成token的函数 一般用于业务登录时生产token给到前端
//@param id int64 用户ID参与生成 token,DecodeToken 函数解密出 id
func CreateToken(id int64) (string, error) {
	tokenData := TokenStruct{
		id,
		time.Now().Unix(),
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 3).Unix(), //token过期时间，不过期设置为0
			Issuer:    "gin-api",
		},
	}
	//
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenData)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

//DecodeToken 解密token的函数
func DecodeToken(token string) (*TokenStruct, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &TokenStruct{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if tokenData, ok := tokenClaims.Claims.(*TokenStruct); ok && tokenClaims.Valid {
			return tokenData, nil
		}
	}

	return nil, err
}

//GenerateToken 生成token
func GenerateToken(userId string) (tokenString string, err error) {
	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Second * 1).Unix(),
		"data": jwt.MapClaims{
			"user_id":     userId,
			"update_time": time.Now().Unix(),
		},
	}

	// 创建一个新的令牌对象，指定签名方法和声明
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密码签名并获得完整的编码令牌作为字符串
	tokenString, err = token.SignedString(jwtSecret)
	return
}

//ParseToken 解析token
func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
