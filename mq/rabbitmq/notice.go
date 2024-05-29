package rabbitmq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logc"
	"io"
	"net/http"
	"strings"
	"time"
)

// WarnRobotNotice 发布到机器人
func WarnRobotNotice(address string, content string) {
	msgBuilder := new(strings.Builder)
	msgBuilder.WriteString("### 企微助手警告，请留意\n")
	msgBuilder.WriteString(fmt.Sprintf("时间：%s\n", time.Now().Format("2006-01-02 15:04:05")))
	msgBuilder.WriteString("内容：")
	msgBuilder.WriteString(content)

	// 构造消息
	notice := struct {
		MsgType  string `json:"msgtype"`
		Markdown struct {
			Content string `json:"content"`
		} `json:"markdown"`
	}{
		MsgType: "markdown",
		Markdown: struct {
			Content string `json:"content"`
		}{
			Content: msgBuilder.String(),
		},
	}

	// json
	jsonData := bytes.Buffer{}
	enc := json.NewEncoder(&jsonData)
	if err := enc.Encode(notice); err != nil {
		logc.Errorf(context.Background(), "编码JSON出错：%s", err.Error())
		return
	}

	go func() {
		req, err := http.NewRequest(http.MethodPost, address, &jsonData)
		if err != nil {
			logc.Errorf(context.Background(), "创建请求出错：%s", err.Error())
			return
		}

		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			logc.Errorf(context.Background(), "发送请求出错：%s", err.Error())
			return
		}
		defer func(Body io.ReadCloser) {
			err = Body.Close()
			if err != nil {
				logc.Errorf(context.Background(), "BodyClose出错：%s", err.Error())
			}
		}(resp.Body)

		respText, err := io.ReadAll(resp.Body)
		if err != nil {
			logc.Errorf(context.Background(), "读取响应体出错：%s", err.Error())
			return
		}

		if resp.StatusCode != http.StatusOK {
			logc.Errorf(context.Background(), "发送告警机器人出错：%s\n", string(respText))
			return
		}

		var respData struct {
			ErrCode int `json:"errcode"`
		}
		if err := json.Unmarshal(respText, &respData); err != nil {
			logc.Errorf(context.Background(), "解析响应体出错：%s", err.Error())
			return
		}
		if respData.ErrCode != 0 {
			logc.Errorf(context.Background(), "发送告警机器人出错：%s\n", string(respText))
			return
		}
	}()
}
