package rmq

import (
	"fmt"
	"github.com/streadway/amqp"
)

var Conn *amqp.Connection

var channelList map[string]*amqp.Channel
var confStr string

func MqConnect(conf string, exchangeList []string) {
	Connect(conf, exchangeList)
}

// conf :   "amqp://admin:你的密码@127.0.0.1:5672/"
func Connect(conf string, exchangeList []string) {
	if confStr == "" {
		confStr = conf
	}
	var err error
	Conn, err = amqp.Dial(confStr)
	if err != nil {
		//fmt.Println("connect rmq", err)
		panic(err)
	}
	if err != nil {
		panic(fmt.Sprintf("connect Channel error: %s", err))
	}
	channelList = make(map[string]*amqp.Channel)
	for _, exchange := range exchangeList {
		channelList[exchange], _ = Conn.Channel()
		err = channelList[exchange].ExchangeDeclare(
			exchange, // name
			"direct", // type
			true,     // durable
			false,    // auto-deleted
			false,    // internal
			false,    // no-wait
			nil,      // arguments
		)
		if err != nil {
			panic(fmt.Sprintf("ch.ExchangeDeclare error: %s", err))
		}
	}

}
func IntoQueue(exchange, routingKey string, message string) bool {
	err := channelList[exchange].Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})

	if err != nil {
		return false
	}

	//fmt.Println(" [x] Sent %s", message)
	return true
}
