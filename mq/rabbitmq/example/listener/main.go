package main

import (
	"flag"
	"fmt"
	"github.com/punpeo/punpeo-lib/mq/rabbitmq"
	"github.com/punpeo/punpeo-lib/mq/rabbitmq/example/listener/config"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
)

var configFile = flag.String("f", "mq/rabbitmq/example/listener/listener.yaml", "Specify the config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)

	listener := rabbitmq.MustNewListener(rabbitmq.RabbitListenerConf{
		RabbitConf:     c.RabbitConf,
		ListenerQueues: c.TestQueues.ListenerQueues,
	}, Handler{})
	serviceGroup := service.NewServiceGroup()
	serviceGroup.Add(listener)
	defer serviceGroup.Stop()
	serviceGroup.Start()
}

type Handler struct {
}

func (h Handler) Consume(message string) error {
	fmt.Printf("listener %s\n", message)
	return nil
}
