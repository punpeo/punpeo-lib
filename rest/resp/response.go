package resp

import (
	jerror2 "github.com/punpeo/punpeo-lib/rest/jerror"
)

// ResponseSuccess is a correct output format.
type ResponseSuccess struct {
	Code uint32      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Success 返回一个正确输出的结构体
func Success(data interface{}) *ResponseSuccess {
	return &ResponseSuccess{jerror2.SuccessCode, "ok", data}
}

// ResponseError is an incorrect output format.
type ResponseError struct {
	Code uint32      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// AppendData 添加错误数据 一般用不到.
func (r *ResponseError) AppendData(data interface{}) *ResponseError {
	r.Data = data
	return r
}

// Error 返回一个错误输出的结构体
func Error(errCode uint32, errMsg string) *ResponseError {
	return &ResponseError{
		Code: errCode,
		Msg:  errMsg,
	}
}
