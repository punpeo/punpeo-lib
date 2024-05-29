package server

import (
	"context"
	"errors"
	pb2 "github.com/punpeo/punpeo-lib/task/gocronnode/pb"
	utils2 "github.com/punpeo/punpeo-lib/task/gocronnode/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Server struct {
	pb2.UnimplementedTaskServer
	taskMap map[string]TaskFunc
}

var keepAlivePolicy = keepalive.EnforcementPolicy{
	MinTime:             10 * time.Second,
	PermitWithoutStream: true,
}

var keepAliveParams = keepalive.ServerParameters{
	MaxConnectionIdle: 30 * time.Second,
	Time:              30 * time.Second,
	Timeout:           3 * time.Second,
}
var Flag = "ZRPC:"

func (s Server) Run(ctx context.Context, req *pb2.TaskRequest) (*pb2.TaskResponse, error) {
	defer func() {
		if err := recover(); err != nil {
			logx.Error(err)
		}
	}()
	logx.Infof("execute cmd start: [id: %d cmd: %s]", req.Id, req.Command)
	var output string
	var err error
	resp := new(pb2.TaskResponse)

	if strings.HasPrefix(req.Command, Flag) {
		var cmdStruct = &utils2.Cmd{}
		cmd := strings.TrimPrefix(strings.TrimSpace(req.Command), Flag)
		cmdStruct, err = utils2.ParseCmd(cmd)
		if err == nil {
			var fn TaskFunc
			fn, err = s.GetTask(cmdStruct.Fn)
			if err == nil {
				f := func() (string, error) {
					return fn(ctx, cmdStruct.Args...)
				}
				output, err = utils2.ExecRpc(ctx, f)
			}
		}
	} else {
		output, err = utils2.ExecShell(ctx, req.Command)
	}
	resp.Output = output
	if err != nil {
		resp.Error = err.Error()
	} else {
		resp.Error = ""
	}
	logx.Infof("execute cmd end: [id: %d cmd: %s err: %s]", req.Id, req.Command, resp.Error)

	return resp, nil
}

func (s Server) Start(addr string) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	opts := []grpc.ServerOption{
		grpc.KeepaliveParams(keepAliveParams),
		grpc.KeepaliveEnforcementPolicy(keepAlivePolicy),
	}

	server := grpc.NewServer(opts...)
	pb2.RegisterTaskServer(server, s)
	reflection.Register(server)
	logx.Infof("server listen on %s", addr)

	go func() {
		err = server.Serve(l)
		if err != nil {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	for {
		s := <-c
		logx.Info("收到信号 -- ", s)
		switch s {
		case syscall.SIGHUP:
			logx.Info("收到终端断开信号, 忽略")
		case syscall.SIGINT, syscall.SIGTERM:
			logx.Info("应用准备退出")
			server.Stop()
			return
		}
	}
}

func (s *Server) RegistTask(tag string, taskFunc TaskFunc) {
	if s.taskMap == nil {
		s.taskMap = make(map[string]TaskFunc)
	}
	if _, ok := s.taskMap[tag]; ok {
		panic("标识位已被使用了")
	}
	s.taskMap[tag] = taskFunc
}

func (s Server) GetTask(tag string) (TaskFunc, error) {
	if _, ok := s.taskMap[tag]; !ok {
		return nil, errors.New("命令未注册")
	}
	return s.taskMap[tag], nil
}
