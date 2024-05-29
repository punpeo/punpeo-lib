package rabbitmq

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"github.com/streadway/amqp"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/queue"
	"log"
	"runtime"
	"strings"
	"time"
)

type (
	ConsumeHandle func(message string) error

	ConsumeHandler interface {
		Consume(message string) error
	}

	RabbitListener struct {
		conn               *amqp.Connection
		channel            *amqp.Channel
		forever            chan bool
		handler            ConsumeHandler
		queues             RabbitListenerConf
		retryDelay         time.Duration
		workerPools        []*ants.Pool
		poolStatsTicker    *time.Ticker
		monitorFrequency   int
		connCloseNotify    chan *amqp.Error
		channelCloseNotify chan *amqp.Error
	}
)

func MustNewListener(listenerConf RabbitListenerConf, handler ConsumeHandler) queue.MessageQueue {
	listener := &RabbitListener{
		queues:           listenerConf,
		handler:          handler,
		forever:          make(chan bool),
		retryDelay:       time.Second * time.Duration(listenerConf.RetryDelay), // 设置重连延迟为5秒
		workerPools:      make([]*ants.Pool, len(listenerConf.ListenerQueues)),
		poolStatsTicker:  nil,
		monitorFrequency: listenerConf.Frequency,
	}

	conn, err := amqp.Dial(getRabbitURL(listenerConf.RabbitConf))
	if err != nil {
		log.Fatalf("failed to connect rabbitmq, error: %+v", err)
	}

	listener.conn = conn
	channel, err := listener.conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %+v", err)
	}
	listener.channel = channel

	name := GetConsumerNames(listenerConf.ListenerQueues)

	listener.connCloseNotify = listener.conn.NotifyClose(make(chan *amqp.Error, 1))
	listener.channelCloseNotify = listener.channel.NotifyClose(make(chan *amqp.Error, 1))

	go listener.checkConnection(name) // 启动连接检查

	return listener
}

func GetConsumerNames(config []ConsumerConf) string {
	names := make([]string, 0)
	for _, conf := range config {
		names = append(names, conf.Name)
	}
	return strings.Join(names, " ")

}

func (q *RabbitListener) Start() {
	q.Consume()
	<-q.forever
}

func panicHandler(err interface{}) {
	// 获取panic发生时的上下文信息
	stackTrace := make([]byte, 4096)
	length := runtime.Stack(stackTrace, false)
	logx.Errorf("ants panic %+v ; stackTrace:%+v", err, string(stackTrace[:length]))
}

func (q *RabbitListener) Consume() {
	// 创建协程池的统计信息定时器
	poolStatsTicker := time.NewTicker(time.Second * time.Duration(q.monitorFrequency))
	q.poolStatsTicker = poolStatsTicker

	for i, que := range q.queues.ListenerQueues {

		err := q.channel.Qos(que.PoolSize, 0, false)
		if err != nil {
			log.Fatalf("failed to set a channel Qos: %+v", err)
		}

		msg, err := q.channel.Consume(
			que.Name,
			"",
			que.AutoAck,
			que.Exclusive,
			que.NoLocal,
			que.NoWait,
			nil,
		)
		if err != nil {
			log.Fatalf("failed to listener, error: %+v", err)
		}

		if que.PoolSize > 1 {
			// 创建协程池
			pool, err := ants.NewPool(que.PoolSize, ants.WithPanicHandler(panicHandler))
			if err != nil {
				log.Fatalf("Failed to create consumer pool: %+v", err)
			}
			q.workerPools[i] = pool

			if que.Monitor {
				go CoroutineHealthCheck(pool, poolStatsTicker, que.Name)
			}

			go func(autoAck bool) {

				defer func() {
					if r := recover(); r != nil {
						logx.Infof("Recovered in handleMessage: %+v", r)
					}
				}()

				for d := range msg {
					err = pool.Submit(handleMessage(d, q, autoAck))
					if err != nil {
						log.Fatalf("failed to Submit, error: %+v", err)
					}
				}

			}(que.AutoAck)
		} else {

			for d := range msg {
				handleMessage(d, q, que.AutoAck)()
			}

		}
	}
}

func CoroutineHealthCheck(pool *ants.Pool, ticker *time.Ticker, name string) {
	defer func() {
		if r := recover(); r != nil {
			// 获取panic发生时的上下文信息
			stackTrace := make([]byte, 4096)
			length := runtime.Stack(stackTrace, false)
			logx.Errorf("coroutine health check is closed: %v,stackTrace:%v,retrying...", r, string(stackTrace[:length]))
			go CoroutineHealthCheck(pool, ticker, name)
		}
	}()
	for range ticker.C {
		// 获取协程池的统计信息
		poolStats := pool.Running()
		isClosed := pool.IsClosed()
		free := pool.Free()
		logx.Infof("queue:%v,running coroutine: %v,free coroutine: %v,coroutine is closed:%v", name, poolStats, free, isClosed)
	}
}

func (q *RabbitListener) checkConnection(name string) {
	defer func() {
		if r := recover(); r != nil {
			// 获取panic发生时的上下文信息
			stackTrace := make([]byte, 4096)
			length := runtime.Stack(stackTrace, false)
			logx.Errorf("rabbitmq health check is closed: %v,stackTrace:%v", r, string(stackTrace[:length]))
			go q.checkConnection(name) // 重新启动连接检查
		}
	}()
	for {
		select {
		case err := <-q.connCloseNotify:
			logx.Infof("queue %v connection closed err:%+v", name, err)
			if len(q.queues.ReportAddress) != 0 {
				WarnRobotNotice(q.queues.ReportAddress, fmt.Sprintf("queue %v connection closed", name))
			}
			q.connRetry(name)
			q.channelRetry(name)
			q.resourceRelease()
			// 重新开始消费消息
			q.Consume()
			logx.Infof("queue %v rabbitmq consume restart", name)
			if len(q.queues.ReportAddress) != 0 {
				WarnRobotNotice(q.queues.ReportAddress, fmt.Sprintf("queue %v rabbitmq consume restart", name))
			}
		case err := <-q.channelCloseNotify:
			logx.Infof("queue %v channel closed err:%+v", name, err)
			if len(q.queues.ReportAddress) != 0 {
				WarnRobotNotice(q.queues.ReportAddress, fmt.Sprintf("queue %v channel closed", name))
			}
			q.channelRetry(name)
			q.resourceRelease()
			// 重新开始消费消息
			q.Consume()
			logx.Infof("queue %v rabbitmq consume restart", name)
			if len(q.queues.ReportAddress) != 0 {
				WarnRobotNotice(q.queues.ReportAddress, fmt.Sprintf("queue %v rabbitmq consume restart", name))
			}
		}
	}
}

func (q *RabbitListener) connRetry(name string) {
	for {
		time.Sleep(q.retryDelay)
		logx.Infof("queue %v connection retrying...", name)
		if len(q.queues.ReportAddress) != 0 {
			WarnRobotNotice(q.queues.ReportAddress, fmt.Sprintf("queue %v connection retrying...", name))
		}
		var err error
		q.conn, err = amqp.Dial(getRabbitURL(q.queues.RabbitConf))
		if err == nil {
			break
		}
		logx.Infof("queue %v failed to reconnect rabbitmq, error: %+v", name, err)
		if len(q.queues.ReportAddress) != 0 {
			WarnRobotNotice(q.queues.ReportAddress, fmt.Sprintf("queue %v failed to reconnect rabbitmq, error: %+v", name, err))
		}
	}
	q.connCloseNotify = q.conn.NotifyClose(make(chan *amqp.Error, 1))
	logx.Infof("queue %v connection retry successful", name)
	if len(q.queues.ReportAddress) != 0 {
		WarnRobotNotice(q.queues.ReportAddress, fmt.Sprintf("queue %v connection retry successful", name))
	}
}

func (q *RabbitListener) channelRetry(name string) {
	for {
		time.Sleep(q.retryDelay)
		logx.Infof("queue %v channel retrying...", name)
		if len(q.queues.ReportAddress) != 0 {
			WarnRobotNotice(q.queues.ReportAddress, fmt.Sprintf("queue %v channel retrying...", name))
		}
		var err error
		q.channel, err = q.conn.Channel()
		if err == nil {
			break
		}
		logx.Infof("queue %v failed to open a channel: %+v", name, err)
		if len(q.queues.ReportAddress) != 0 {
			WarnRobotNotice(q.queues.ReportAddress, fmt.Sprintf("queue %v failed to open a channel: %+v", name, err))
		}
	}
	q.channelCloseNotify = q.channel.NotifyClose(make(chan *amqp.Error, 1))
	logx.Infof("queue %v channel retry successful", name)
	if len(q.queues.ReportAddress) != 0 {
		WarnRobotNotice(q.queues.ReportAddress, fmt.Sprintf("queue %v channel retry successful", name))
	}
}

func (q *RabbitListener) resourceRelease() {
	// 关闭协程池和协程池的统计信息定时器
	for _, pool := range q.workerPools {
		if pool != nil && !pool.IsClosed() {
			pool.Release()
		}
	}

	if q.poolStatsTicker != nil {
		q.poolStatsTicker.Stop()
	}

	q.workerPools = make([]*ants.Pool, len(q.queues.ListenerQueues))
	q.poolStatsTicker = nil
}

type taskFunc func()

func handleMessage(d amqp.Delivery, q *RabbitListener, autoAck bool) taskFunc {
	return func() {
		if err := q.handler.Consume(string(d.Body)); err != nil {
			logx.Errorf("Error on consuming: %s, error: %+v", string(d.Body), err)
			if !autoAck {
				err = d.Nack(false, true)
				if err != nil {
					logx.Errorf("Error to nack: %s, error: %+v", string(d.Body), err)
				}
			}

		} else {
			if !autoAck {
				err = d.Ack(false)
				if err != nil {
					logx.Errorf("Error to ack: %s, error: %+v", string(d.Body), err)
				}
			}

		}
	}

}

func (q *RabbitListener) Stop() {
	q.channel.Close()
	q.conn.Close()
	// 关闭协程池和相关统计
	for _, pool := range q.workerPools {
		if pool != nil && !pool.IsClosed() {
			pool.Release()
		}

	}
	if q.poolStatsTicker != nil {
		q.poolStatsTicker.Stop()
	}

	close(q.forever)
}
