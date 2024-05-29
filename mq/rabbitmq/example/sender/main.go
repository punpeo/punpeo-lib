package main

import (
	"encoding/json"
	"github.com/punpeo/punpeo-lib/mq/rabbitmq"
)

func main() {
	conf := rabbitmq.RabbitConf{
		Username: "guest",
		Password: "guest",
		Host:     "127.0.0.1",
		Port:     5672,
	}

	queue := rabbitmq.SenderConfig{
		ExchangeName: "qwzs_queue_test",
		QueueName:    "qwzs_queue_test",
		RouterKey:    "qwzs_queue_test",
		ContentType:  "application/json",
	}

	data := map[string]interface{}{
		"msg": "json test 111",
	}

	msg, _ := json.Marshal(data)
	err := rabbitmq.Publish(msg, conf, queue)
	if err != nil {
		return
	}

}
