package restyclient

import "time"

type Option func(o *options)

type options struct {
	retryCount       int           // 重试的次数
	retryWaitTime    time.Duration // 重试等待时间，默认3秒
	retryMaxWaitTime time.Duration // 重试最大等待时间 10秒
	timeout          time.Duration // 默认超时时间，最大3秒
}

func RetryCount(retryCount int) Option {
	return func(o *options) {
		o.retryCount = retryCount
	}
}

func RetryWaitTime(retryWaitTime time.Duration) Option {
	return func(o *options) {
		o.retryWaitTime = retryWaitTime
	}
}

func RetryMaxWaitTime(retryMaxWaitTime time.Duration) Option {
	return func(o *options) {
		o.retryMaxWaitTime = retryMaxWaitTime
	}
}

func Timeout(timeout time.Duration) Option {
	return func(o *options) {
		o.timeout = timeout
	}
}
