package xerr

var message map[uint32]string

func init() {
	message = map[uint32]string{
		OK:                            "SUCCESS",
		SERVER_COMMON_ERROR:           "服务器开小差啦,稍后再来试一试",
		REUQEST_PARAM_ERROR:           "参数错误",
		LOGIN_EXPIRE_ERROR:            "登录信息过期，请重新登录",
		SIGN_ERROR:                    "签名错误",
		MISSED_DATA_PERMISSIONS_ERROR: "没有数据权限",
		HTTP_SERVER_REQUEST_ERROR:     "http服务请求错误",
		THIRD_SERVER_REQUEST_ERROR:    "第三方服务错误",
		MISSED_FUNC_PERMISSIONS_ERROR: "没有功能权限",
	}

}

func MapErrMsg(errcode uint32) string {
	if msg, ok := message[errcode]; ok {
		return msg
	} else {
		return "服务器开小差啦,稍后再来试一试"
	}
}

func IsCodeErr(errcode uint32) bool {
	if _, ok := message[errcode]; ok {
		return true
	} else {
		return false
	}
}
