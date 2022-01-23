package e

const (
	SUCCESS      = 200
	ERROR        = 500
	PARAM_ERROR  = 400
	NOT_FOUND    = 404
	LOGIC_ERROR  = 600
	AUTH_FAIL    = 401
	TOKEN_EXPIRE = 10002
	SIGN_ERROR = 10003
)

var message = map[int]string{
	SUCCESS:      "ok",
	ERROR:        "fail",
	PARAM_ERROR:  "请求参数错误", //各业务逻辑各自的参数错误，可以重写该message
	NOT_FOUND:    "404 error",
	LOGIC_ERROR:  "业务逻辑错误", //各业务逻辑各自的错误，可以重写该message
	AUTH_FAIL:    "鉴权失败",
	TOKEN_EXPIRE: "token过期",
	SIGN_ERROR: "签名错误",
}

func GetMessage(code int) string {
	msg, ok := message[code]
	if ok {
		return msg
	}

	return message[ERROR]
}
