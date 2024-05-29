package main

import (
	"context"
	"fmt"
	server2 "github.com/punpeo/punpeo-lib/task/gocronnode/server"
	"time"
)

func main() {
	s := new(server2.Server)
	s.RegistTask("test", func(ctx context.Context, params ...string) (string, error) {
		for _, v := range params {
			fmt.Println(v)
		}
		return "ok", nil
	})
	s.RegistTask("test2", func(ctx context.Context, params ...string) (string, error) {
		myTimer := time.NewTicker(time.Second) // 启动定时器
		for {
			select {
			case <-ctx.Done():
				myTimer.Stop()
				fmt.Println("stop")
				return "ok", nil
			case a := <-myTimer.C:
				fmt.Println(a.String())
			}
		}
	})
	s.Start("0.0.0.0:8080")
}
