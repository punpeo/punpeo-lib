package requests

import (
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var jsonContentType = "application/json; charset=utf-8"

func TestParseForm(t *testing.T) {
	var v struct {
		Name    string  `form:"name"`
		Age     int     `form:"age"`
		Percent float64 `form:"percent,optional"`
	}

	r, err := http.NewRequest(http.MethodGet, "/a?name=hello&age=18&percent=3.4", nil)
	assert.Nil(t, err)
	assert.Nil(t, ParseAndValidate(r, &v))
	assert.Equal(t, "hello", v.Name)
	assert.Equal(t, 18, v.Age)
	assert.Equal(t, 3.4, v.Percent)
}

func TestParseForm_Error(t *testing.T) {
	var v struct {
		Name    string  `form:"name"`
		Age     int     `form:"age"`
		Percent float64 `form:"percent,optional"`
	}

	r, err := http.NewRequest(http.MethodGet, "/a?percent=3.4", nil)
	assert.Nil(t, err)
	err = ParseAndValidate(r, &v)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "字段name不能为空")

	// 类型不匹配
	r, err = http.NewRequest(http.MethodGet, "/a?name=aa&age=aaa&percent=3.4", nil)
	err = ParseAndValidate(r, &v)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "aaa类型异常")
}

func TestParseJsonBody(t *testing.T) {
	t.Run("validate test", func(t *testing.T) {
		// 没有validate
		var v struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		body := `{"name":"aa", "age": 5}`
		r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		r.Header.Set(httpx.ContentType, jsonContentType)

		assert.Nil(t, ParseAndValidate(r, &v))
		assert.Equal(t, "aa", v.Name)
		assert.Equal(t, 5, v.Age)
	})

	t.Run("validate test", func(t *testing.T) {
		// 有validate
		var v struct {
			Name string `json:"name" validate:"min=3,max=10"`
			Age  int    `json:"age" validate:"gte=10,lte=100"`
		}

		body := `{"name":"aa", "age": 5}`
		r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		r.Header.Set(httpx.ContentType, jsonContentType)

		err := ParseAndValidate(r, &v)
		assert.Error(t, err)
		t.Log(err.Error())
	})

}
