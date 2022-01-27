package request

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	zhongwen "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	errors "github.com/pkg/errors"
)

var v *validator.Validate
var trans ut.Translator

func init() {
	// 中文翻译
	zh := zhongwen.New()
	uni := ut.New(zh, zh)
	trans, _ = uni.GetTranslator("zh")

	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		// 验证器注册翻译器
		zh_translations.RegisterDefaultTranslations(v, trans)
		// 自定义验证方法
		v.RegisterValidation("checkMobile", checkMobile)
		v.RegisterTranslation("checkMobile", trans, checkMobileMsg, checkMobileRegister)
	}
}

func ParseRequest(c *gin.Context, request interface{}) error {
	err := c.ShouldBind(request)
	var errStr string

	if err != nil {
		switch err.(type) {
		case validator.ValidationErrors:
			errStr = Translate(err.(validator.ValidationErrors))
		case *json.UnmarshalTypeError:
			unmarshalTypeError := err.(*json.UnmarshalTypeError)
			errStr = fmt.Errorf("%s 类型错误，期望类型 %s", unmarshalTypeError.Field, unmarshalTypeError.Type.String()).Error()
		default:
			errStr = errors.New("unknown error.").Error()
		}

		return errors.New(errStr)
	}

	return nil
}

func Translate(errs validator.ValidationErrors) string {
	var errList []string
	for _, e := range errs {
		// can translate each error one at a time.
		errList = append(errList, e.Translate(trans))
	}

	//不全部返回错误，返回第一个错误
	return errList[0]

	//return strings.Join(errList, "|")
}

func checkMobile(fl validator.FieldLevel) bool {
	mobile := strconv.Itoa(int(fl.Field().Uint()))
	re := `^1[3456789]\d{9}$`
	r := regexp.MustCompile(re)
	return r.MatchString(mobile)
}

func checkMobileMsg(ut ut.Translator) error {
	return ut.Add("checkMobile", "{0}长度不等于11位或{1}格式错误!", true) // see universal-translator for details
}

func checkMobileRegister(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T("checkMobile", fe.Field(), fe.Field())

	return t
}
