package createqueue

import (
	"encoding/json"
	"github.com/punpeo/punpeo-lib/mq/rabbitmq"
	"testing"
)

func Test_CreateRoutingQueue(t *testing.T) {
	conf := rabbitmq.RabbitConf{
		Username: "guest",
		Password: "guest",
		Host:     "127.0.0.1",
		Port:     5672,
		//VHost:    "/",
	}
	queue := rabbitmq.SenderConfig{
		ExchangeName: "qwzs_queue_direct01",
		QueueName:    "qwzs_queue_direct01",
		RouterKey:    "qwzs_queue_direct01",
		ContentType:  "application/json",
	}
	err1 := rabbitmq.CreateRoutingQueue(conf, queue)
	if err1 != nil {
		t.Errorf("create queue fail: %v", err1)
	}

	data := map[string]interface{}{
		"msg": "json test 111",
	}

	msg, _ := json.Marshal(data)
	err := rabbitmq.Publish(msg, conf, queue)
	if err != nil {
		t.Errorf("send msg fail: %v", err)
	}
}
