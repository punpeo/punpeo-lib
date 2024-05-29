package httpxpro

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

// GetFormValue 从form中获取指定 key 参数
func GetFormValue(r *http.Request, key string) string {
	form, err := httpx.GetFormValues(r)
	if err != nil {
		return ""
	}
	value, ok := form[key]
	if !ok {
		return ""
	}
	res, ok := value.(string)
	if !ok {
		return ""
	}
	return res
}
