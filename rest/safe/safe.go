package safe

import (
	"github.com/zeromicro/go-zero/core/logx"
	"runtime"
)

func Go(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != any(nil) {
				buf := make([]byte, 64<<10)
				buf = buf[:runtime.Stack(buf, false)]
				logx.Errorf("Panic Recover(%s) %s", err, buf)
			}
		}()
		f()
	}()

}
