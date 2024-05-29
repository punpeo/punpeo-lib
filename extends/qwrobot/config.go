package qwrobot

import (
	"errors"
	"github.com/zeromicro/go-zero/core/collection"
	"time"
)

var (
	// ErrFrequency 超过报警发送频率限制.
	ErrFrequency = errors.New("too high frequency")
	// DefaultFrequencyConf 默认发送频率限制 5分钟10次
	DefaultFrequencyConf = NewFrequencyConf(10, 60*5)
)

type FrequencyConf struct {
	// Frequency 窗口时间发送数限制
	Frequency int
	// Duration 窗口时间 单位S
	Duration int
	rw       *collection.RollingWindow
}

// NewFrequencyConf 发送频率限制
// 单位时间 duration 秒内限制最大发送 frequency 条
func NewFrequencyConf(frequency, duration int) *FrequencyConf {
	if duration <= 0 {
		panic("duration time must be greater than 0")
	}
	frequencyConf := &FrequencyConf{
		Frequency: frequency,
		Duration:  duration,
	}
	// 窗口大小为1S
	frequencyConf.rw = collection.NewRollingWindow(duration, time.Second)
	return frequencyConf
}
