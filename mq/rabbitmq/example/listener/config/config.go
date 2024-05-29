package config

import "github.com/punpeo/punpeo-lib/mq/rabbitmq"

type Config struct {
	TestQueues rabbitmq.ConsumerConfig
	RabbitConf rabbitmq.RabbitConf
}
