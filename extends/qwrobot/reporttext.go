package qwrobot

import (
	"github.com/zeromicro/go-zero/core/collection"
)

// ReportMessage 发送一个文本报错
func ReportMessage(url, message string) error {
	robot := Robot{
		Url:     url,
		Msgtype: "text",
		Text: Text{
			Content: message,
		},
	}
	return robot.SendMessageWithText()
}

// ReportMessageFrequency 发送一个文本报错并且限制报警频率
func ReportMessageFrequency(url, message string, frequencyConf *FrequencyConf) error {
	var sum float64
	frequencyConf.rw.Reduce(func(b *collection.Bucket) {
		sum += b.Sum
	})
	if int(sum) >= frequencyConf.Frequency {
		return ErrFrequency
	}

	robot := Robot{
		Url:     url,
		Msgtype: "text",
		Text: Text{
			Content: message,
		},
	}
	frequencyConf.rw.Add(1)
	return robot.SendMessageWithText()
}
