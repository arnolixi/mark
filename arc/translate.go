package arc

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

// InitTrans 初始化翻译器
var trans ut.Translator

func init() {
	err := NewITrans("zh")
	if err != nil {
		panic(errors.Wrap(err, "初始化校验翻译器失败!"))
	}
}

func NewITrans(locale string) (err error) {
	// 修改gin框架中的Validator属性，实现自定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册一个获取json tag的自定义方法
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		zhT := zh.New() // 中文翻译器
		enT := en.New() // 英文翻译器

		uni := ut.New(enT, zhT, enT)
		trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed\n", locale)
		}

		//_ = v.RegisterTranslation("required_with", trans, func(ut ut.Translator) error {
		//	return ut.Add("required_with", "{0} 为必填字段!", true)
		//}, func(ut ut.Translator, fe validator.FieldError) string {
		//	t, _ := ut.T("required_with", fe.Field())
		//	return t
		//})
		//_ = v.RegisterTranslation("required_without", trans, func(ut ut.Translator) error {
		//	return ut.Add("required_without", "{0} 为必填字段!", true)
		//}, func(ut ut.Translator, fe validator.FieldError) string {
		//	t, _ := ut.T("required_without", fe.Field())
		//	return t
		//})
		//_ = v.RegisterTranslation("required_without_all", trans, func(ut ut.Translator) error {
		//	return ut.Add("required_without_all", "{0} 为必填字段!", true)
		//}, func(ut ut.Translator, fe validator.FieldError) string {
		//	t, _ := ut.T("required_without_all", fe.Field())
		//	return t
		//})

		switch locale {
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(v, trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		}
		return
	}
	return
}

func T() ut.Translator {
	return trans
}
