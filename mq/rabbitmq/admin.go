package rabbitmq

import (
	"github.com/streadway/amqp"
	"github.com/zeromicro/go-zero/core/logx"
)

type Admin struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func MustNewAdmin(rabbitMqConf RabbitConf) *Admin {
	var admin Admin
	conn, err := amqp.Dial(getRabbitURL(rabbitMqConf))
	if err != nil {
		logx.Errorf("failed to connect rabbitmq, error: %v", err)
		return nil
	}

	admin.conn = conn
	channel, err := admin.conn.Channel()
	if err != nil {
		logx.Errorf("failed to open a channel, error: %v", err)
		conn.Close()
		return nil
	}

	admin.channel = channel
	return &admin
}

func (q *Admin) DeclareExchange(conf ExchangeConf, args amqp.Table) error {
	return q.channel.ExchangeDeclare(
		conf.ExchangeName,
		conf.Type,
		conf.Durable,
		conf.AutoDelete,
		conf.Internal,
		conf.NoWait,
		args,
	)
}

func (q *Admin) DeclareQueue(conf QueueConf, args amqp.Table) error {
	_, err := q.channel.QueueDeclare(
		conf.Name,
		conf.Durable,
		conf.AutoDelete,
		conf.Exclusive,
		conf.NoWait,
		args,
	)

	return err
}

func (q *Admin) Bind(queueName string, routekey string, exchange string, notWait bool, args amqp.Table) error {
	return q.channel.QueueBind(
		queueName,
		routekey,
		exchange,
		notWait,
		args,
	)
}

// 创建routring模式的队列，仅执行一次做队列初始化
func CreateRoutingQueue(rabbitMqConf RabbitConf, smqConf SenderConfig) error {
	admin := MustNewAdmin(rabbitMqConf)
	exchangeConf := ExchangeConf{
		ExchangeName: smqConf.ExchangeName,
		Type:         "direct",
		Durable:      true,
		AutoDelete:   false,
		Internal:     false,
		NoWait:       false,
	}
	err := admin.DeclareExchange(exchangeConf, nil)
	if err != nil {
		logx.Errorf("failed to declare exchange, error: %v", err)
		return err
	}
	queueConf := QueueConf{
		Name:       smqConf.QueueName,
		Durable:    true,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
	}
	err = admin.DeclareQueue(queueConf, nil)
	if err != nil {
		logx.Errorf("failed to declare queue, error: %v", err)
		return err
	}
	err = admin.Bind(smqConf.QueueName, smqConf.RouterKey, smqConf.ExchangeName, false, nil)
	if err != nil {
		logx.Errorf("failed to bind routing , error: %v", err)
		return err
	}
	return nil
}
