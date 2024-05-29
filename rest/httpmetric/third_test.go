package httpmetirc_test

import (
	httpmetirc "github.com/punpeo/punpeo-lib/rest/httpmetric"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/prometheus"
	"github.com/zeromicro/go-zero/core/timex"
	"runtime"
	"testing"
	"time"
)

func TestInc(t *testing.T) {
	startAgent()
	httpmetirc.Inc("/user/name", 200, 0, "GET")
	httpmetirc.Inc("/user/data", 200, 0, "POST")
}

func TestAdd(t *testing.T) {
	startAgent()
	startTime := timex.Now()
	httpmetirc.Observe(startTime, "/user/name", "POST")
	httpmetirc.Observe(startTime, "/user/data", "POST")
}

func TestManyRequest(t *testing.T) {
	startAgent()
	go func() {
		for i := 0; i < 1000; i++ {
			go func() {
				httpmetirc.Inc("/user/name", 200, 0, "GET")

			}()
			time.Sleep(time.Second)
		}
	}()

	for {
		logx.Infof("goroutine num is %d", runtime.NumGoroutine())
		time.Sleep(time.Second)
	}
}

func startAgent() {
	prometheus.StartAgent(prometheus.Config{
		Host: "0.0.0.0",
		Port: 9090,
		Path: "/metrics",
	})
}
