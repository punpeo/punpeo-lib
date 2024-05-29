package httpmetirc

import (
	"github.com/zeromicro/go-zero/core/metric"
)

const (
	serverNamespace = "third_http"
)

var (
	metricThirdWxReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "third http client requests duration(ms).",
		Labels:    []string{"srv", "path", "method"}, // 服务名，请求路径，请求方式
		Buckets:   []float64{100, 200, 300, 500, 1000, 2000, 3000, 5000},
	})

	metricThirdWxReqCodeTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "third http client requests code count.",
		Labels:    []string{"srv", "path", "http_code", "srv_code", "method"}, // 服务名，请求路径，http_code , srv_code码，请求方式
	})
)

type ThirdResp struct {
	Code      int         `json:"code,omitempty"`
	Msg       string      `json:"msg,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	RequestId string      `json:"request_id,omitempty"`
}
