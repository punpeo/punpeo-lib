package httpmetirc

import (
	"github.com/zeromicro/go-zero/core/timex"
	"net/url"
	"strconv"
	"time"
)

const (
	TuoHeng = "TuoHeng"
)

func Observe(startTime time.Duration, path, method string) {
	parsedURL, _ := url.Parse(path)
	metricThirdWxReqDur.Observe(timex.Since(startTime).Milliseconds(), TuoHeng, parsedURL.Path, method)
}

func Inc(path string, httpCode int, srvCode int, method string) {
	parsedURL, _ := url.Parse(path)
	metricThirdWxReqCodeTotal.Inc(TuoHeng, parsedURL.Path, strconv.Itoa(httpCode), strconv.Itoa(srvCode), method)
}
