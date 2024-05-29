//go:build !windows
// +build !windows

package utils

import (
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/net/context"
	"os/exec"
	"strings"
	"syscall"
)

type Result struct {
	output string
	err    error
}

// 执行shell命令，可设置执行超时时间
func ExecShell(ctx context.Context, command string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	resultChan := make(chan Result)
	go func() {
		output, err := cmd.CombinedOutput()
		resultChan <- Result{string(output), err}
	}()
	select {
	case <-ctx.Done():
		if cmd.Process.Pid > 0 {
			syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		}
		return "", errors.New("timeout killed")
	case result := <-resultChan:
		return result.output, result.err
	}
}

func ConvertEncoding(outputGBK string) string {
	// windows平台编码为gbk，需转换为utf8才能入库
	outputUTF8, ok := GBK2UTF8(outputGBK)
	if ok {
		return outputUTF8
	}

	return outputGBK
}

type Cmd struct {
	Fn   string
	Args []string
}

func ExecRpc(ctx context.Context, fn func() (string, error)) (string, error) {
	var resultChan chan Result = make(chan Result)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logx.Error(err)
			}
		}()
		output, err := fn()
		resultChan <- Result{string(output), err}
	}()
	select {
	case <-ctx.Done():
		return "", errors.New("timeout killed")
	case result := <-resultChan:
		return ConvertEncoding(result.output), result.err
	}
}

func ParseCmd(cmd string) (*Cmd, error) {
	flags := strings.Split(cmd, " ")
	if len(flags) == 0 {
		return nil, errors.New("参数错误")
	}
	var args []string

	//将字符串切片追加到os.Args中
	for key, value := range flags {
		if key == 0 {
			continue
		}
		if strings.HasPrefix(value, "-") {
			if strings.Contains(value, "=") {
				s := strings.Split(value, "=")
				if len(s) > 0 {
					args = append(args, s[len(s)-1])
				}
			}
		} else {
			args = append(args, value)
		}
	}
	return &Cmd{
		Fn:   flags[0],
		Args: args,
	}, nil
}
