package render

import (
	jerror2 "github.com/punpeo/punpeo-lib/rest/jerror"
	"github.com/punpeo/punpeo-lib/rest/resp"
	"github.com/punpeo/punpeo-lib/zeropro/response"
	"net/http"
)

// JsonpResult 输出jsonp
func JsonpResult(r *http.Request, w http.ResponseWriter, resp interface{}, err error) {
	if err == nil {
		res := resp.Success(resp)
		response.WriteJsonp(r, w, http.StatusOK, res)
	} else {
		switch e := err.(type) {
		case *jerror2.CodeError:
			response.WriteJsonp(r, w, http.StatusOK, resp.Error(e.Code, e.Msg))
		default:
			response.WriteJsonp(r, w, http.StatusOK, resp.Error(jerror2.DefaultCode, e.Error()))
		}

	}
}

// JsonpErrorResult 输出错误.
func JsonpErrorResult(r *http.Request, w http.ResponseWriter, err error, data ...interface{}) {
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
	response.WriteJsonp(r, w, http.StatusBadRequest, apiErr)
}
