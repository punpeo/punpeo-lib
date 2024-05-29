package restyclient

import (
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"
)

type zHttp struct {
	cli  *resty.Client
	opts options
}

var (
	HdrUserAgentKey       = http.CanonicalHeaderKey("User-Agent")
	HdrAcceptKey          = http.CanonicalHeaderKey("Accept")
	HdrContentTypeKey     = http.CanonicalHeaderKey("Content-Type")
	HdrContentLengthKey   = http.CanonicalHeaderKey("Content-Length")
	HdrContentEncodingKey = http.CanonicalHeaderKey("Content-Encoding")
	HdrLocationKey        = http.CanonicalHeaderKey("Location")
)

const (
	FormContentType = "application/x-www-form-urlencoded"
	JsonContentType = "application/json"
	PlainTextType   = "text/plain; charset=utf-8"
)

func NewClient(cli *resty.Client, opts ...Option) *resty.Client {
	o := options{
		retryCount:       0,
		retryWaitTime:    3 * time.Second,
		retryMaxWaitTime: 10 * time.Second,
		timeout:          3 * time.Second,
	}

	for _, opt := range opts {
		opt(&o)
	}

	zh := &zHttp{
		cli:  cli,
		opts: o,
	}
	return zh.client()
}

func NewJsonContent(cli *resty.Client, opts ...Option) *resty.Client {
	client := NewClient(cli, opts...)
	client.SetHeader(HdrContentTypeKey, JsonContentType)
	return client
}

func (z *zHttp) client() *resty.Client {
	if z.opts.retryCount != 0 {
		z.cli.SetRetryCount(z.opts.retryCount)

		// 设置重试等待时间
		if z.opts.retryWaitTime != 0 {
			z.cli.SetRetryWaitTime(z.opts.retryWaitTime)
		}
		if z.opts.retryMaxWaitTime != 0 {
			z.cli.SetRetryMaxWaitTime(z.opts.retryMaxWaitTime)
		}
	}

	// 设置超时控制
	if z.opts.timeout != 0 {
		z.cli.SetTimeout(z.opts.timeout)
	}
	return z.cli
}
