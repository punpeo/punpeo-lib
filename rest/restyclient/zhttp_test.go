package restyclient

import (
	"github.com/go-resty/resty/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"runtime"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {

	go func() {
		for i := 0; i < 10; i++ {
			client := resty.New()
			resp, err := NewClient(client).R().Get("https://baidu.com/")
			if err != nil {
				panic(any(err))
			}
			logx.Infof("%d", resp.StatusCode())

			client.GetClient().CloseIdleConnections()
			time.Sleep(1 * time.Second)
		}
	}()

	for {
		logx.Info("goroutines num is ", runtime.NumGoroutine())
		time.Sleep(1 * time.Second)
	}
}

func TestTimeout(t *testing.T) {
	resp, err := NewClient(resty.New(), Timeout(100*time.Millisecond)).R().Get("https://httpbin.org/get")
	if err != nil {
		logx.Info(err)
	}
	logx.Infof("%s", resp.Body())
}

func TestTryAgain(t *testing.T) {
	resp, err := NewJsonContent(resty.New(), Timeout(5*time.Second),
		RetryCount(3)).R().Get("https://httpbin.org/get")
	if err != nil {
		logx.Info(err)
	}
	logx.Infof("%s", resp.Body())
}
