package safe

import (
	"github.com/zeromicro/go-zero/core/logx"
	"testing"
)

func TestGo(t *testing.T) {
	ch := make(chan struct{})
	Go(func() {
		logx.Info("do something.....")
		defer func() { ch <- struct{}{} }()
		panic(any("协程报错。。。"))
	})

	<-ch
	logx.Info("recover data.")
}
