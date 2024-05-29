package restyclient

import (
	"github.com/punpeo/punpeo-lib/rest/xerr"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"
	"time"
)

func TestHttpConnectClose(t *testing.T) {
	go func() {
		for i := 0; i < 10; i++ {
			getReq("https://www.baidu.com/")
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for i := 0; i < 10; i++ {
			getReq("https://www.qq.com/")
			time.Sleep(1 * time.Second)
		}
	}()

	for i := 0; i < 30; i++ {
		logx.Infof("goroutines num is %d", runtime.NumGoroutine())
		time.Sleep(1 * time.Second)
	}
}

// mockServer 创建一个模拟服务器用于测试
func mockServer() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/test-form-data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":1000,"msg":"success","data":{"foo":"bar"}}`))
	})

	return httptest.NewServer(handler)
}

// TestHttpPostSendFormData 测试 HttpPostSendFormData 方法
func TestHttpPostSendFormData(t *testing.T) {
	server := mockServer()
	defer server.Close()

	// 设置测试数据
	url := server.URL + "/test-form-data"
	params := map[string]string{"field1": "value1", "field2": "value2"}
	timeout := 5 * time.Second

	// 调用方法
	result, err := HttpPostSendFormData(url, params, timeout)

	// 验证结果
	if err != nil {
		t.Errorf("HttpPostSendFormData returned an error: %v", err)
	}
	if uint32(result.Code) != xerr.OK {
		t.Errorf("Expected code %d, got %d", xerr.OK, result.Code)
	}
	if result.Msg != "success" {
		t.Errorf("Expected message 'success', got '%s'", result.Msg)
	}
}

func getReq(url string) {
	_, err := HttpGetSendQuery(url, nil, 15*time.Second)
	if err != nil {
		panic(any(err))
	}
}
