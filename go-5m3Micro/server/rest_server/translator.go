package rest_server

import (
	"reflect"
	"strings"

	"go-5m3Micro/pkg/errors"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTrans "github.com/go-playground/validator/v10/translations/en"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
)

func (srv *Server) initTranslator(local string) error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		var err error
		// 注册 自定义 获取 Json 字段名 方法
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
		// 中文翻译器
		zhT := zh.New()
		// 英文翻译器
		enT := en.New()
		// 第一个参数为 备用翻译器
		uni := ut.New(zhT, enT, zhT)
		var exits bool
		if srv.trans, exits = uni.GetTranslator(local); exits {
			switch local {
			case "zh":
				err = zhTrans.RegisterDefaultTranslations(v, srv.trans)
				if err != nil {
					return errors.Errorf("init zh translation error: %v", err)
				}
			case "en":
				err = enTrans.RegisterDefaultTranslations(v, srv.trans)
				if err != nil {
					return errors.Errorf("init en translation error: %v", err)
				}
			default:
				err = zhTrans.RegisterDefaultTranslations(v, srv.trans)
				if err != nil {
					return errors.Errorf("init zh translation error: %v", err)
				}
			}
		} else {
			return errors.Errorf("get translator error")
		}
	} else {
		return errors.Errorf("init translator error")
	}
	return nil
}
