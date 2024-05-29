package render

import (
	jerror2 "github.com/punpeo/punpeo-lib/rest/jerror"
	"github.com/punpeo/punpeo-lib/rest/resp"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

// JsonResult 输出jsonp
func JsonResult(r *http.Request, w http.ResponseWriter, resp interface{}, err error) {
	if err == nil {
		res := resp.Success(resp)
		httpx.WriteJson(w, http.StatusOK, res)
	} else {
		switch e := err.(type) {
		case *jerror2.CodeError:
			httpx.WriteJson(w, http.StatusOK, resp.Error(e.Code, e.Msg))
		default:
			httpx.WriteJson(w, http.StatusOK, resp.Error(jerror2.DefaultCode, e.Error()))
		}

	}
}

// JsonErrorResult 输出错误.
func JsonErrorResult(r *http.Request, w http.ResponseWriter, err error, data ...interface{}) {
	var apiErr *resp.ResponseError
	switch e := err.(type) {
	case *jerror2.CodeError:
		apiErr = resp.Error(e.Code, err.Error())
	default:
		apiErr = resp.Error(jerror2.ParamCode, err.Error())
	}
	if len(data) == 1 {
		apiErr.AppendData(data[0])
	}
	httpx.WriteJson(w, http.StatusBadRequest, apiErr)
}
