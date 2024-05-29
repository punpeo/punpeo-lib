package response

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type message struct {
	Name string `json:"name"`
}

func TestWriteJsonp(t *testing.T) {
	w := tracedResponseWriter{
		headers: make(map[string][]string),
	}
	msg := message{Name: "anyone"}
	r := httptest.NewRequest(http.MethodGet, "/a?callback=aa", nil)
	WriteJsonp(r, &w, http.StatusOK, msg)
	assert.Equal(t, w.builder.String(), `aa({"name":"anyone"});`)

	w2 := tracedResponseWriter{
		headers: make(map[string][]string),
	}
	msg2 := message{Name: "anyone"}
	r2 := httptest.NewRequest(http.MethodGet, "/a", nil)
	WriteJsonp(r2, &w2, http.StatusOK, msg2)
	assert.Equal(t, w2.builder.String(), `{"name":"anyone"}`)
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
