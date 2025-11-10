package validation

import (
	"regexp"
	
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func RegisterMobile(translator ut.Translator) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", validateMobile)
		_ = v.RegisterTranslation("mobile", translator, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 非法的手机号", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}
}

func validateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	if ok, _ := regexp.MatchString(`^1\d{10}$`, mobile); ok {
		return ok
	}
	return false
}
