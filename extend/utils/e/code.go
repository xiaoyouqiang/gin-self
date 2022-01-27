package e

const (
	Success      = 200
	Error        = 500
	ParamError  = 400
	NotFound    = 404
	LogicError  = 600
	AuthFail    = 401
	TokenExpire = 10002
	SignError = 10003
)

var message = map[int]string {
	Success:      "ok",
	Error:        "fail",
	ParamError:  "请求参数错误", //各业务逻辑各自的参数错误，可以重写该message
	NotFound:    "404 error",
	LogicError:  "业务逻辑错误", //各业务逻辑各自的错误，可以重写该message
	AuthFail:    "鉴权失败",
	TokenExpire: "token过期",
	SignError:   "签名错误",
}

func GetMessage(code int) string {
	msg, ok := message[code]
	if ok {
		return msg
	}

	return message[Error]
}
