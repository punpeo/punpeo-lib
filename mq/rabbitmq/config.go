package rabbitmq

import "fmt"

type RabbitConf struct {
	Username      string
	Password      string
	Host          string
	Port          int
	VHost         string `json:",optional"`
	Frequency     int    `json:",default=5"`
	RetryDelay    int    `json:",default=3"` // 队列重连间隔
	ReportAddress string `json:",optional"`  // 队列故障通知 - 企业微信机器人地址
}

type RabbitListenerConf struct {
	RabbitConf
	ListenerQueues []ConsumerConf
}

type ConsumerConfig struct {
	ListenerQueues []ConsumerConf
}

type SenderConfig struct {
	ExchangeName string
	QueueName    string
	RouterKey    string
	ContentType  string
}

type ConsumerConf struct {
	Name      string
	AutoAck   bool `json:",default=true"`
	Exclusive bool `json:",default=false"`
	// Set to true, which means that messages sent by producers in the same connection
	// cannot be delivered to consumers in this connection.
	NoLocal bool `json:",default=false"`
	// Whether to block processing
	NoWait   bool `json:",default=false"`
	PoolSize int  `json:",default=10"`    // 协程池大小
	Monitor  bool `json:",default=false"` // 是否开启协程池信息统计
}

type RabbitSenderConf struct {
	RabbitConf
	ContentType string `json:",default=text/plain"` // MIME content type
}

type QueueConf struct {
	Name       string
	Durable    bool `json:",default=true"`
	AutoDelete bool `json:",default=false"`
	Exclusive  bool `json:",default=false"`
	NoWait     bool `json:",default=false"`
}

type ExchangeConf struct {
	ExchangeName string
	Type         string `json:",options=direct|fanout|topic|headers"` // exchange type
	Durable      bool   `json:",default=true"`
	AutoDelete   bool   `json:",default=false"`
	Internal     bool   `json:",default=false"`
	NoWait       bool   `json:",default=false"`
	Queues       []QueueConf
}

func getRabbitURL(rabbitConf RabbitConf) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/%s", rabbitConf.Username, rabbitConf.Password,
		rabbitConf.Host, rabbitConf.Port, rabbitConf.VHost)
}
