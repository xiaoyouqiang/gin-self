package request_validator

type Register struct {
	Mobile   uint   `form:"mobile" binding:"required,checkMobile"` //checkMobile自定义验证函数 request.checkMobile
	Password string `form:"password" binding:"required,gte=6"`
	Email    string `json:"email" binding:"required,email"`
}
