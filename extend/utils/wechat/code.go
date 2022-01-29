package wechat

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"gin-self/conf"
	"gin-self/extend/self_redis"
	"gin-self/extend/utils/request"
	"regexp"
	"strings"
	"time"
)

const (
	codeToSession = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=%s"
)

var (
	appId     string
	appSecret string
	grantType string
)

type CodeSessionResponse struct {
	ErrCode int64 `json:"errcode"`
	ErrMsg string `json:"errmsg"`
	UnionId string `json:"unionid"`
	OpenId string `json:"openid"`
	SessionKey string `json:"session_key"`
}

func init()  {
	appId = conf.Wechat["app_id"]
	appSecret = conf.Wechat["secret"]
	grantType = "authorization_code"
}

//GetKeyInfoByCode 通过 code获取 openid unionid session_key 信息
func GetKeyInfoByCode(ctx context.Context, code string) (CodeSessionResponse,error) {

	return CodeSessionResponse{
		ErrCode: 0,
		ErrMsg: "",
		UnionId: "F5g1bVAABfgu9jWsDiyHLJZJyM1",
		OpenId: "obmig4qV9wMuZuSbYoQNGDeFinKg1",
		SessionKey: "tiihtNczf5v6AKRyjwEUhQ==",
	},nil

	var result CodeSessionResponse

	data,err := self_redis.GetConn("master").Get("ye_ma_wx_code_" + code)
	if err != nil {
		return CodeSessionResponse{},err
	} else if data != "" {
		err = json.Unmarshal([]byte(data), &result)
		if err != nil {
			return CodeSessionResponse{},err
		}
		return result,nil
	}

	url := fmt.Sprintf(codeToSession,appId,appSecret,code,grantType)

	resp,err := request.HttpGetByCtx(ctx,url)

	if err != nil {
		return CodeSessionResponse{},err
	}

	err = json.Unmarshal([]byte(resp), &result)
	if err != nil {
		return CodeSessionResponse{},err
	}

	if result.ErrCode != 0 {
		return CodeSessionResponse{},errors.New(result.ErrMsg)
	}

	self_redis.GetConn("master").Set("ye_ma_wx_code_" + code, resp, time.Duration(10) * time.Minute)

	return result,nil
}

//DecryptUserInfo 解码微信个人资料 需要返回 JSON 数据类型时 使用 true, 需要返回 map 数据类型时 使用 false
func DecryptUserInfo(encryptedData,iv,sessionKey string, isJSON bool) (interface{}, error) {
	sessionKey = strings.Replace(strings.TrimSpace(sessionKey), " ", "+", -1)
	if len(sessionKey) != 24 {
		return nil, errors.New("sessionKey length is error")
	}
	aesKey, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		return nil, errors.New("decodeBase64Error")
	}
	iv = strings.Replace(strings.TrimSpace(iv), " ", "+", -1)
	if len(iv) != 24 {
		return nil, errors.New("iv length is error")
	}
	aesIv, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, errors.New("decodeBase64Error")
	}
	encryptedData = strings.Replace(strings.TrimSpace(encryptedData), " ", "+", -1)
	aesCipherText, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, errors.New("decodeBase64Error")
	}
	aesPlantText := make([]byte, len(aesCipherText))

	aesBlock, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, errors.New("illegalBuffer")
	}

	mode := cipher.NewCBCDecrypter(aesBlock, aesIv)
	mode.CryptBlocks(aesPlantText, aesCipherText)
	aesPlantText = PKCS7UnPadding(aesPlantText)

	var decrypted map[string]interface{}

	re := regexp.MustCompile(`[^\{]*(\{.*\})[^\}]*`)
	aesPlantText = []byte(re.ReplaceAllString(string(aesPlantText), "$1"))
	err = json.Unmarshal(aesPlantText, &decrypted)
	if err != nil {
		return nil, errors.New("decodeJsonError")
	}

	//if decrypted["watermark"].(map[string]interface{})["appid"] != wxCrypt.AppId {
	//	return nil, errors.New("appId is not match")
	//}

	if isJSON {
		return string(aesPlantText), nil
	}

	return decrypted, nil
}

// PKCS7UnPadding return unpadding []Byte plantText
func PKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	if length > 0 {
		unPadding := int(plantText[length-1])
		return plantText[:(length - unPadding)]
	}
	return plantText
}