package validate

import (
	"errors"
	"fmt"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

// jzValidator is a simple validator for Jz
// 底层通过 https://github.com/go-playground/validator 实现数据验证
var jzValidator JzValidator

type JzValidator struct {
	once          sync.Once
	validate      *validator.Validate
	validateTrans ut.Translator
}

func (v *JzValidator) ValidateStruct(obj interface{}) error {
	if obj == nil {
		return nil
	}
	v.lazyInit()
	return v.validate.Struct(obj)
}

func (v *JzValidator) lazyInit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validateTrans = validateTrans(v.validate)
		v.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			// 获取tag中的参数名称 友好输出
			name := strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
			if name == "" {
				name = strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			}
			if name == "-" {
				return ""
			}
			return name
		})
	})
}

// validateTrans 数据验证翻译器.
func validateTrans(validate *validator.Validate) ut.Translator {
	uni := ut.New(zh.New())
	// 翻译器
	trans, _ := uni.GetTranslator("zh")
	err := zhTranslations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		fmt.Println(err)
	}
	return trans
}

// JzValidate 校验数据.
func JzValidate(obj interface{}) error {
	return jzValidator.ValidateStruct(obj)
}

// JzValidateDealErr 对验证数据错误进行翻译以及相关数据处理
func JzValidateDealErr(err error) error {
	jzValidator.lazyInit()
	validateErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		// 非验证错误
		if _, ok := err.(*strconv.NumError); ok {
			// 类似这种错误
		}
		return errors.New("参数异常-未知异常")
	}

	for _, value := range validateErrs.Translate(jzValidator.validateTrans) {
		return errors.New(value)
	}
	return nil
}
