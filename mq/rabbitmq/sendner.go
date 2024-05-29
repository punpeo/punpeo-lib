package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"github.com/zeromicro/go-zero/core/logx"
)

type Sender struct {
	conn        *amqp.Connection
	channel     *amqp.Channel
	ContentType string
}

func NewSender(rabbitMqConf RabbitSenderConf) (*Sender, error) {
	sender := &Sender{ContentType: rabbitMqConf.ContentType}
	conn, err := amqp.Dial(getRabbitURL(rabbitMqConf.RabbitConf))
	if err != nil {
		logx.Error(fmt.Sprintf("rabbitSender_conn_err : %v", err))
		return sender, err
	}
	sender.conn = conn

	channel, err := conn.Channel()
	if err != nil {
		logx.Error(fmt.Sprintf("rabbitSender_channel_err : %v", err))
		sender.Close()
		return sender, err
	}
	sender.channel = channel

	return sender, nil
}

func (q *Sender) Send(exchange string, routeKey string, msg []byte, delay ...int64) error {
	headers := amqp.Table{}
	if len(delay) > 0 {
		headers["x-delay"] = delay[0]
	}
	return q.channel.Publish(
		exchange,
		routeKey,
		false,
		false,
		amqp.Publishing{
			Headers:     headers,
			ContentType: q.ContentType,
			Body:        msg,
		},
	)
}

func (q *Sender) Close() {
	if q.channel != nil {
		q.channel.Close()
	}
	if q.conn != nil {
		q.conn.Close()
	}
}

func Publish(body []byte, rabbitConf RabbitConf, senderQueue SenderConfig, delay ...int64) error {
	defer func() {
		if err := recover(); err != nil {
			logx.Info(fmt.Sprintf(" rabbitSender_publish_err body:%v,err:%v", body, err))
		}
	}()
	conf := RabbitSenderConf{RabbitConf: rabbitConf, ContentType: senderQueue.ContentType}
	sender, err := NewSender(conf)
	if err != nil {
		return err
	}
	defer sender.Close()

	err = sender.Send(senderQueue.ExchangeName, senderQueue.RouterKey, body, delay...)
	return err
}
