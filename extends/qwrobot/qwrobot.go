package qwrobot

import (
	"fmt"
	"github.com/punpeo/punpeo-lib/rest/restyclient"
	"github.com/punpeo/punpeo-lib/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

type Robot struct {
	Url      string
	Msgtype  string
	MarkDown MarkDown
	Text     Text
}

// markdown类型
type MarkDown struct {
	Content string
}

// text类型
type Text struct {
	Content             string
	MentionedList       []string
	MentionedMobileList []string
}

var contentMap = make(map[int64]string)

const RobotContentSysErrorType = 1
const RobotContentSysErrorMsg = "<font color=\"warning\">系统异常警告</font>\n " +
	">请求Url：<font color=\"comment\">%+v</font>\n " +
	">请求Body：<font color=\"comment\">%+v</font>\n " +
	">异常内容：%+v"

func init() {
	contentMap[RobotContentSysErrorType] = RobotContentSysErrorMsg
}

// 发送Markdown消息
func (r *Robot) SendMessageWithMarkDown() error {
	body := map[string]interface{}{
		"msgtype": utils.Ternary(r.Msgtype == "", "markdown", r.Msgtype),
		"markdown": map[string]interface{}{
			"content": r.MarkDown.Content,
		},
	}
	_, err := restyclient.HttpPostSendJson(r.Url, body, 0)
	if err != nil {
		logx.Error(fmt.Sprintf("SendMessageWithMarkDown err : %v", err))
		return nil
	}

	return nil
}

// 发送Text消息
func (r *Robot) SendMessageWithText() error {
	body := map[string]interface{}{
		"msgtype": utils.Ternary(r.Msgtype == "", "text", r.Msgtype),
		"text": map[string]interface{}{
			"content":               r.Text.Content,
			"mentioned_list":        r.Text.MentionedList,
			"mentioned_mobile_list": r.Text.MentionedMobileList,
		},
	}
	_, err := restyclient.HttpPostSendJson(r.Url, body, 0)
	if err != nil {
		logx.Error(fmt.Sprintf("SendMessageWithText err : %v", err))
		return nil
	}

	return nil
}

func (r *Robot) GetContentByType(messageType int64) string {
	content := "未知错误"
	if v, ok := contentMap[messageType]; ok {
		content = v
	}

	return content
}
