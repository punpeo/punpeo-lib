package render

import (
	"errors"
	jerror2 "github.com/punpeo/punpeo-lib/rest/jerror"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// message 表示测试输出的内容.
type message struct {
	Name string `json:"name"`
}

var rWithCallback = httptest.NewRequest(http.MethodGet, "/a?callback=aa", nil)
var rNormal = httptest.NewRequest(http.MethodGet, "/a", nil)

func TestJsonResult(t *testing.T) {
	tests := []struct {
		name      string
		request   *http.Request
		data      interface{} // 输出数据
		assertMsg string      // 断言
	}{
		{
			name:      "输出结构体数据",
			request:   rNormal,
			data:      message{Name: "Jz"},
			assertMsg: `{"code":1000,"msg":"ok","data":{"name":"Jz"}}`,
		},
		{
			name:      "输出空数据",
			request:   rNormal,
			data:      nil,
			assertMsg: `{"code":1000,"msg":"ok","data":null}`,
		},
		{
			name:      "输出map",
			request:   rNormal,
			data:      map[string]interface{}{"name": "jz", "age": 18},
			assertMsg: `{"code":1000,"msg":"ok","data":{"age":18,"name":"jz"}}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := tracedResponseWriter{
				headers: make(map[string][]string),
			}
			JsonResult(test.request, &w, test.data, nil)
			assert.Equal(t, http.StatusOK, w.code)
			t.Log(test.assertMsg)
			t.Log(w.builder.String())
			assert.Equal(t, test.assertMsg, w.builder.String())
		})
	}
}

func TestJsonErrorResult(t *testing.T) {
	tests := []struct {
		name      string
		request   *http.Request
		assertMsg string // 断言
		err       error
		data      interface{}
	}{
		{
			name:      "默认参数错误",
			request:   rNormal,
			assertMsg: `{"code":1002,"msg":"参数异常"}`,
			err:       jerror2.ParamErr,
		},
		{
			name:      "自定义错误",
			request:   rNormal,
			assertMsg: `{"code":1002,"msg":"错误文案"}`,
			err:       errors.New("错误文案"),
		},
		{
			name:      "参数错误 增加数据",
			request:   rNormal,
			assertMsg: `{"code":1002,"msg":"错误文案","data":{"name":"Jz"}}`,
			err:       errors.New("错误文案"),
			data:      message{Name: "Jz"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := tracedResponseWriter{
				headers: make(map[string][]string),
			}
			JsonErrorResult(test.request, &w, test.err, test.data)
			t.Log(test.assertMsg)
			t.Log(w.builder.String())
			assert.Equal(t, test.assertMsg, w.builder.String())
		})
	}
}

type tracedResponseWriter struct {
	headers     map[string][]string
	builder     strings.Builder
	hasBody     bool
	code        int
	lessWritten bool
	wroteHeader bool
	err         error
}

func (w *tracedResponseWriter) Header() http.Header {
	return w.headers
}

func (w *tracedResponseWriter) Write(bytes []byte) (n int, err error) {
	if w.err != nil {
		return 0, w.err
	}

	n, err = w.builder.Write(bytes)
	if w.lessWritten {
		n -= 1
	}
	w.hasBody = true

	return
}

func (w *tracedResponseWriter) WriteHeader(code int) {
	if w.wroteHeader {
		return
	}
	w.wroteHeader = true
	w.code = code
}
