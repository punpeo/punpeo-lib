package response

import (
	"encoding/json"
	"github.com/punpeo/punpeo-lib/utils/conv"
	"github.com/punpeo/punpeo-lib/zeropro/httpxpro"
	"github.com/zeromicro/go-zero/rest/httpx"
	"html/template"
	"net/http"
)

const (
	JsonpContentType = "application/javascript; charset=utf-8"
)

// WriteJsonp writes v as jsonp string into w with code.输出jsonp.
func WriteJsonp(r *http.Request, w http.ResponseWriter, code int, v interface{}) {
	callback := httpxpro.GetFormValue(r, "callback")
	if callback == "" {
		httpx.WriteJson(w, http.StatusOK, v)
		return
	}
	bs, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(httpx.ContentType, JsonpContentType)
	w.WriteHeader(code)

	callback = template.JSEscapeString(callback)

	if _, err = w.Write(conv.StringToBytes(callback)); err != nil {
		panic(err)
	}

	if _, err = w.Write(conv.StringToBytes("(")); err != nil {
		panic(err)
	}

	if _, err := w.Write(bs); err != nil {
		panic(err)
	}

	if _, err = w.Write(conv.StringToBytes(");")); err != nil {
		panic(err)
	}
}
