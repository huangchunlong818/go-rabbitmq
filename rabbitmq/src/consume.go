package src

import (
	"context"
	"github.com/avast/retry-go"
	errorx "github.com/huangchunlong818/go-rabbitmq/rabbitmq/error"
	"strconv"
	"time"
)

var (
	// 默认消费参数
	defaultConsumeArgs = consumeArgs{
		qosPrefetchCount:   30,
		qosPrefetchSize:    0,
		qosGlobal:          false,
		autoAck:            false,
		ackMultiple:        false,
		nackMultiple:       false,
		nackRequeue:        false,
		queueBindNoWait:    false,
		queueDurable:       true,
		exchangeDurable:    true,
		exchangeAutoDelete: false,
		exchangeInternal:   false,
		exchangeNoWait:     false,
		consumeExclusive:   false,
		consumeNoLocal:     false,
		consumeNoWait:      false,
	}
)

// consumeArgs 消费参数
type consumeArgs struct {
	qosPrefetchCount   int  // 每次从队列中获取的消息数量
	qosPrefetchSize    int  // 每次从队列中获取的消息大小
	qosGlobal          bool // 对当前mq连接所有channel生效
	autoAck            bool // 是否自动确认消息
	ackMultiple        bool // 是否批量确认消息
	nackMultiple       bool // 是否批量拒绝当前和之后的消息
	nackRequeue        bool // 重新消费策略
	queueBindNoWait    bool // 是否阻塞等待队列创建完成
	queueDurable       bool // 队列是否持久化
	exchangeDurable    bool // 交换机是否持久化
	exchangeAutoDelete bool // 交换机是否自动删除
	exchangeInternal   bool // 交换机是否内部使用
	exchangeNoWait     bool // 是否阻塞等待交换机创建完成
	consumeExclusive   bool // 是否独占消费
	consumeNoLocal     bool // 是否消费本地消息
	consumeNoWait      bool // 是否阻塞等待消费者创建完成

	// 以下参数不默认
	exchangeName string // 交换机名称
	exchangeType string // 交换机类型
	queueName    string // 队列名称
	routingKey   string // 路由键
}

type ConsumeOption func(*consumeArgs)

// WithConsumeQosPrefetchCount 设置每次从队列中获取的消息数量
func WithConsumeQosPrefetchCount(count int) ConsumeOption {
	return func(args *consumeArgs) {
		args.qosPrefetchCount = count
	}
}

// WithConsumeQosPrefetchSize 设置每次从队列中获取的消息大小
func WithConsumeQosPrefetchSize(size int) ConsumeOption {
	return func(args *consumeArgs) {
		args.qosPrefetchSize = size
	}
}

// WithConsumeQosGlobal 设置对当前mq连接所有channel生效
func WithConsumeQosGlobal(global bool) ConsumeOption {
	return func(args *consumeArgs) {
		args.qosGlobal = global
	}
}

// WithConsumeAutoAck 设置是否自动确认消息
func WithConsumeAutoAck(autoAck bool) ConsumeOption {
	return func(args *consumeArgs) {
		args.autoAck = autoAck
	}
}

// WithConsumeAckMultiple 设置是否批量确认消息
func WithConsumeAckMultiple(ackMultiple bool) ConsumeOption {
	return func(args *consumeArgs) {
		args.ackMultiple = ackMultiple
	}
}

// WithConsumeNackMultiple 设置是否批量拒绝当前和之后的消息
func WithConsumeNackMultiple(nackMultiple bool) ConsumeOption {
	return func(args *consumeArgs) {
		args.nackMultiple = nackMultiple
	}
}

// WithConsumeNackRequeue 设置重新消费策略
func WithConsumeNackRequeue(nackRequeue bool) ConsumeOption {
	return func(args *consumeArgs) {
		args.nackRequeue = nackRequeue
	}
}

// WithConsumeQueueBindNoWait 设置是否阻塞等待队列创建完成
func WithConsumeQueueBindNoWait(queueBindNoWait bool) ConsumeOption {
	return func(args *consumeArgs) {
		args.queueBindNoWait = queueBindNoWait
	}
}

// WithConsumeExchangeDurable 设置交换机是否持久化
func WithConsumeExchangeDurable(exchangeDurable bool) ConsumeOption {
	return func(args *consumeArgs) {
		args.exchangeDurable = exchangeDurable
	}
}

// WithConsumeExchangeAutoDelete 设置交换机是否自动删除
func WithConsumeExchangeAutoDelete(exchangeAutoDelete bool) ConsumeOption {
	return func(args *consumeArgs) {
		args.exchangeAutoDelete = exchangeAutoDelete
	}
}

// WithConsumeExchangeInternal 设置交换机是否内部使用
func WithConsumeExchangeInternal(exchangeInternal bool) ConsumeOption {
	return func(args *consumeArgs) {
		args.exchangeInternal = exchangeInternal
	}
}

// WithConsumeExchangeNoWait 设置是否阻塞等待交换机创建完成
func WithConsumeExchangeNoWait(exchangeNoWait bool) ConsumeOption {
	return func(args *consumeArgs) {
		args.exchangeNoWait = exchangeNoWait
	}
}

// WithConsumeExclusive 设置是否独占消费
func WithConsumeExclusive(consumeExclusive bool) ConsumeOption {
	return func(args *consumeArgs) {
		args.consumeExclusive = consumeExclusive
	}
}

// WithConsumeNoLocal 设置是否消费本地消息
func WithConsumeNoLocal(consumeNoLocal bool) ConsumeOption {
	return func(args *consumeArgs) {
		args.consumeNoLocal = consumeNoLocal
	}
}

// WithConsumeNoWait 设置是否阻塞等待消费者创建完成
func WithConsumeNoWait(consumeNoWait bool) ConsumeOption {
	return func(args *consumeArgs) {
		args.consumeNoWait = consumeNoWait
	}
}

// WithConsumeQueueDurable 设置队列是否持久化
func WithConsumeQueueDurable(queueDurable bool) ConsumeOption {
	return func(args *consumeArgs) {
		args.queueDurable = queueDurable
	}
}

// dataHandleFunc 消息处理函数
type dataHandleFunc func([]byte) error

// ConsumeTopic 启动Topic消费者
func (c *RabmqConnPool) ConsumeTopic(exchangeName, queueName, routingKey string, handler dataHandleFunc, num int, timeInterval time.Duration, ops ...ConsumeOption) {
	c.consumeWithExchangeType(exchangeName, "topic", queueName, routingKey, handler, num, timeInterval, ops...)
}

// ConsumeFanout 启动Fanout消费者
func (c *RabmqConnPool) ConsumeFanout(exchangeName, queueName string, handler dataHandleFunc, num int, timeInterval time.Duration, ops ...ConsumeOption) {
	c.consumeWithExchangeType(exchangeName, "fanout", queueName, "", handler, num, timeInterval, ops...)
}

// ConsumeDirect 启动Direct消费者
func (c *RabmqConnPool) ConsumeDirect(exchangeName, queueName, routingKey string, handler dataHandleFunc, num int, timeInterval time.Duration, ops ...ConsumeOption) {
	c.consumeWithExchangeType(exchangeName, "direct", queueName, routingKey, handler, num, timeInterval, ops...)
}

// ConsumeDirect 启动Direct消费者
func (c *RabmqConnPool) consumeWithExchangeType(exchangeName, exchangeType, queueName, routingKey string, handler dataHandleFunc, num int, timeInterval time.Duration, ops ...ConsumeOption) {
	ctx := context.Background()
	consu := defaultConsumeArgs
	for _, op := range ops {
		op(&consu)
	}
	consu.exchangeName = exchangeName
	consu.exchangeType = exchangeType
	consu.queueName = queueName
	consu.routingKey = routingKey

	for {
		connWrapper, err := c.getConnWrapper(ctx)
		if err != nil {
			errorx.ErrorPush("getConnWrapper failed：" + err.Error())
			<-time.After(RETRY_GET_CONNDURATION)
			continue
		}
		err = c.startConsumeWithConn(ctx, connWrapper, &consu, handler, num, timeInterval)
		if err != nil {
			errorx.ErrorPush("startConsumeWithChannel failed：" + err.Error())
			<-time.After(RETRY_GET_CONNDURATION)
		}
	}
}

// startConsumeWithChannel 通过connWrapper启动消费者
func (c *RabmqConnPool) startConsumeWithConn(ctx context.Context, connWrapper *mqConnWrapper, args *consumeArgs, handler dataHandleFunc, num int, timeInterval time.Duration) error {
	channel, err := connWrapper.conn.Channel()
	if err != nil {
		return err
	}
	// 如果是临时队列，需要在消费者启动时创建队列
	if !args.queueDurable {
		queue, err := channel.QueueDeclare(args.queueName, false, true, true, false, nil)
		if err != nil {
			return err
		}
		args.queueName = queue.Name
	}
	// 交换机定义
	err = channel.ExchangeDeclare(args.exchangeName, args.exchangeType, args.exchangeDurable, args.exchangeAutoDelete, args.exchangeInternal, args.exchangeNoWait, nil)
	if err != nil {
		return err
	}
	// 设置每次从队列中获取的消息数量
	err = channel.Qos(args.qosPrefetchCount, args.qosPrefetchSize, args.qosGlobal)
	if err != nil {
		return err
	}
	//绑定队列到 exchange 中
	err = channel.QueueBind(args.queueName, args.routingKey, args.exchangeName, args.queueBindNoWait, nil)
	if err != nil {
		return err
	}
	//消费消息
	messages, err := channel.Consume(args.queueName, args.routingKey, args.autoAck, args.consumeExclusive, args.consumeNoLocal, args.consumeNoWait, nil)
	if err != nil {
		return err
	}
	// 消费消息
	for d := range messages {
		//获取当前队列最大的消费次数
		var retryCount uint
		if num == 0 {
			retryCount = 1 //默认只消费一次
		}

		// 使用 retry.Do 封装处理逻辑
		errs := retry.Do(
			func() error {
				return handler(d.Body)
			},
			retry.Attempts(retryCount), // 总共执行几次
			retry.Delay(timeInterval),  // 重试间隔
			retry.OnRetry(func(n uint, err error) {
				errorx.ErrorPush("消费消息失败：queueName=" + args.queueName + ",msg=" + string(d.Body) + "，error：" + err.Error() + "，当前消费次数：" + strconv.Itoa(int(n+1)))
			}),
		)
		//ret := handler(d.Body)
		if !args.autoAck {
			if errs == nil {
				d.Ack(args.ackMultiple) //参数是否批量确认
			} else {
				d.Nack(args.nackMultiple, args.nackRequeue) //第一个参数是否批量，第二个参数是否重新入队
			}
		}
	}
	return nil
}
