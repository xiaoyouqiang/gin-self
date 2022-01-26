package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"gin-self/conf"
	"gin-self/extend/utils/request"
)

const (
	codeToSession = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=%s"
)

var (
	appId     string
	appSecret string
	grantType string
)

type codeSessionResponse struct {
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
func GetKeyInfoByCode(ctx context.Context, code string) (codeSessionResponse,error) {

	url := fmt.Sprintf(codeToSession,appId,appSecret,code,grantType)

	resp,err := request.HttpGetByCtx(ctx,url)

	if err != nil {
		return codeSessionResponse{},err
	}

	var result codeSessionResponse
	json.Unmarshal([]byte(resp), &result)

	return result,nil
}