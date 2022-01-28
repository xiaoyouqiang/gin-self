package wechat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gin-self/conf"
	"gin-self/extend/self_redis"
	"gin-self/extend/utils/request"
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

	//return codeSessionResponse{
	//	ErrCode: 0,
	//	ErrMsg: "",
	//	UnionId: "F5g1bVAABfgu9jWsDiyHLJZJyM",
	//	OpenId: "obmig4qV9wMuZuSbYoQNGDeFinKg",
	//	SessionKey: "",
	//},nil

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

//DecodeEncryptedData 解码授权数据 获取手机号
func DecodeEncryptedData(encryptedData, iv, sessionKey string) (string,error) {

	return "",nil
}