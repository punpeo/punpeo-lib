package restyclient

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/punpeo/punpeo-lib/rest/xerr"
	"net/http"
	"time"
)

type HttpResultMap struct {
	Code      int64                  `json:"code"`
	Msg       string                 `json:"msg"`
	Timestamp int64                  `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

type HttpResultList struct {
	Code      int64                    `json:"code"`
	Msg       string                   `json:"msg"`
	Timestamp int64                    `json:"timestamp"`
	Data      []map[string]interface{} `json:"data"`
}

type HttpResult struct {
	Code      int64       `json:"code"`
	Msg       string      `json:"msg"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// HttpPostSendJson 目前链接没办法复用，每个请求都需要create and close
func HttpPostSendJson(url string, params map[string]interface{}, timeout time.Duration) (HttpResult, error) {
	httpResult := HttpResult{}

	jsonBody, err := json.Marshal(params)
	if err != nil {
		//logx.Logger.Error(xerr.MapErrMsg(xerr.REUQEST_PARAM_ERROR), err)
		return httpResult, err
	}

	// 发起请求
	client := resty.New()
	// 修复链接释放慢问题，高并发场景会出现协程挂起，造成系统毛刺过高
	defer client.GetClient().CloseIdleConnections()
	resp, err := client.SetTimeout(timeout).R().
		SetHeader("Content-Type", "application/json").
		SetBody(string(jsonBody)).
		SetResult(&httpResult).
		ForceContentType("application/json").
		Post(url)

	if err != nil {
		//logx.Logger.Error(xerr.MapErrMsg(xerr.HTTP_SERVER_NET_ERROR), err)
		return httpResult, err
	}
	if resp.StatusCode() != 200 {
		//logx.Logger.Error(xerr.MapErrMsg(xerr.HTTP_SERVER_NET_ERROR), err)
		return httpResult, err
	}

	if httpResult.Code != 1000 {
		//logx.Logger.Error(xerr.MapErrMsg(xerr.HTTP_SERVER_REQUEST_ERROR), err)
		return httpResult, err
	}

	return httpResult, nil
}

// HttpGetSendQuery  目前链接没办法复用，每个请求都需要create and close
func HttpGetSendQuery(url string, params map[string]string, timeout time.Duration) (HttpResult, error) {
	httpResult := HttpResult{}

	// 发起请求
	client := resty.New()
	// 修复链接释放慢问题，高并发场景会出现协程挂起，造成系统毛刺过高
	defer client.GetClient().CloseIdleConnections()
	resp, err := client.SetTimeout(timeout).R().
		SetQueryParams(params).
		SetResult(&httpResult).
		Get(url)

	if err != nil {
		//logx.Logger.Error(xerr.MapErrMsg(xerr.HTTP_SERVER_NET_ERROR), err)
		return httpResult, err
	}
	if resp.StatusCode() != 200 {
		//logx.Logger.Error(xerr.MapErrMsg(xerr.HTTP_SERVER_NET_ERROR), err)
		return httpResult, err
	}

	if httpResult.Code != 1000 {
		//logx.Logger.Error(xerr.MapErrMsg(xerr.HTTP_SERVER_REQUEST_ERROR), err)
		return httpResult, err
	}

	return httpResult, nil
}

// HttpPostSendFormData 发送 multipart/form-data 类型的 POST 请求
func HttpPostSendFormData(url string, params map[string]string, timeout time.Duration) (HttpResult, error) {
	httpResult := HttpResult{}

	// 创建 resty 客户端
	client := resty.New()
	// 修复链接释放慢问题
	defer client.GetClient().CloseIdleConnections()

	//跳过证书验证,本地服务测试
	//tr := &http.Transport{
	//	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	//}
	//client.SetTransport(tr)

	// 构建请求
	req := client.SetTimeout(timeout).R().SetResult(&httpResult)
	// 设置不需要验证 https 证书

	// 添加表单数据
	for key, value := range params {
		req.SetFormData(map[string]string{key: value})
	}

	// 发送请求
	resp, err := req.Post(url)

	// 错误处理
	if err != nil {
		return httpResult, err
	}
	if resp.StatusCode() != http.StatusOK {
		return httpResult, err
	}
	if uint32(httpResult.Code) != xerr.OK {
		return httpResult, err
	}

	return httpResult, nil
}
