package jerror

type CodeError struct {
	Code uint32 `json:"code"`
	Msg  string `json:"msg"`
}

func NewCodeError(code uint32, msg string) error {
	return &CodeError{Code: code, Msg: msg}
}

func (e *CodeError) Error() string {
	return e.Msg
}
