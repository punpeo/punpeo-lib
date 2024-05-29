package jerror

const (
	SuccessCode = 1000
	DefaultCode = 1001 // 默认错误
	ParamCode   = 1002 // 参数异常
	SignCode    = 1003 // 签名异常
)

var (
	DefaultErr = NewDefaultError("发生错误")
	ParamErr   = NewCodeError(ParamCode, "参数异常")
	SignErr    = NewCodeError(SignCode, "签名错误")
)

// NewDefaultError 返回一个默认错误 错误码 DefaultCode
func NewDefaultError(msg string) error {
	return NewCodeError(DefaultCode, msg)
}

// NewParamError 返回一个参数异常错误
func NewParamError(msg string) error {
	return NewCodeError(DefaultCode, msg)
}
