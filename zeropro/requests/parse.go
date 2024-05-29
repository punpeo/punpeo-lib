package requests

import (
	"fmt"
	"github.com/punpeo/punpeo-lib/utils/validate"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"regexp"
)

var regNotSet = regexp.MustCompile(`field (\w+) is not set`)                 // 参数不存在
var regParsedErr = regexp.MustCompile(`the value "(\w+)" cannot parsed as`)  // form类型不匹配
var regMismatch = regexp.MustCompile(`error: type mismatch for field (\w+)`) // json类型不匹配

// ParseAndValidate parses the request and Validate parameters. 对go-zero中的Parse 在验证方面做了增强.
func ParseAndValidate(r *http.Request, v interface{}) error {
	if err := httpx.ParsePath(r, v); err != nil {
		return err
	}

	if err := httpx.ParseForm(r, v); err != nil {
		return transParseErr(err)
	}

	if err := httpx.ParseHeaders(r, v); err != nil {
		return err
	}

	if err := httpx.ParseJsonBody(r, v); err != nil {
		return transParseErr(err)
	}

	// 验证
	if err := validate.JzValidate(v); err != nil {
		return validate.JzValidateDealErr(err)
	}

	return nil
}

func transParseErr(err error) error {
	msg := err.Error()

	// 没有找到
	findField := regNotSet.FindStringSubmatch(msg)
	if len(findField) == 2 {
		return fmt.Errorf("字段%s不能为空", findField[1])
	}
	// parseErr
	findField2 := regParsedErr.FindStringSubmatch(msg)
	if len(findField2) == 2 {
		return fmt.Errorf("%s类型异常", findField2[1])
	}
	findField3 := regMismatch.FindStringSubmatch(msg)
	if len(findField3) == 2 {
		return fmt.Errorf("字段%s类型异常", findField3[1])
	}

	// 待补充
	return err
}
