package src

import (
	"context"
	"errors"
	errorx "github.com/huangchunlong818/go-rabbitmq/rabbitmq/error"
	"github.com/streadway/amqp"
)

var (
	defaultProducerArgs = producerArgs{
		isConfirm:           true,
		publishMandatory:    false,
		publishImmediate:    false,
		publishDeliveryMode: amqp.Persistent,
		publishPriority:     0,
	}
)

type producerArgs struct {
	isConfirm           bool  // 确认模式
	publishMandatory    bool  // 未匹配到队列处理方法
	publishImmediate    bool  // 路由失败时处理方法
	publishDeliveryMode uint8 // 消息传递失败后是否重试投递
	publishPriority     uint8 // 消息优先级
	delay               int   // 消息延迟时间（单位：毫秒） **新增延迟选项**
	// 以下参数不默认
	exchangeName string // 交换机名称
	routingKey   string // 路由键
}

type ProducerOption func(*producerArgs)

// WithProducerDelay 设置消息延迟 // **新增延迟选项**
func WithProducerDelay(delay int) ProducerOption {
	return func(args *producerArgs) {
		args.delay = delay
	}
}

// WithProducerIsConfirm 设置确认模式
func WithProducerIsConfirm(isConfirm bool) ProducerOption {
	return func(args *producerArgs) {
		args.isConfirm = isConfirm
	}
}

// WithProducerPublishMandatory 设置未匹配到队列处理方法
func WithProducerPublishMandatory(publishMandatory bool) ProducerOption {
	return func(args *producerArgs) {
		args.publishMandatory = publishMandatory
	}
}

// WithProducerPublishImmediate 设置路由失败时处理方法
func WithProducerPublishImmediate(publishImmediate bool) ProducerOption {
	return func(args *producerArgs) {
		args.publishImmediate = publishImmediate
	}
}

// WithProducerPublishDeliveryMode 设置消息传递失败后是否重试投递
func WithProducerPublishDeliveryMode(publishDeliveryMode uint8) ProducerOption {
	return func(args *producerArgs) {
		args.publishDeliveryMode = publishDeliveryMode
	}
}

// WithProducerPublishPriority 设置消息优先级
func WithProducerPublishPriority(publishPriority uint8) ProducerOption {
	return func(args *producerArgs) {
		args.publishPriority = publishPriority
	}
}

// ProducerTopic  Topic消息生产
func (c *RabmqConnPool) ProducerTopic(ctx context.Context, exchangeName, routingKey string, data []byte, ops ...ProducerOption) error {
	return c.producerPublish(ctx, exchangeName, routingKey, data, ops...)
}

// ProducerDirect  Direct消息生产
func (c *RabmqConnPool) ProducerDirect(ctx context.Context, exchangeName, routingKey string, data []byte, ops ...ProducerOption) error {
	return c.producerPublish(ctx, exchangeName, routingKey, data, ops...)
}

// ProducerFanout  Fanout消息生产
func (c *RabmqConnPool) ProducerFanout(ctx context.Context, exchangeName string, data []byte, ops ...ProducerOption) error {
	return c.producerPublish(ctx, exchangeName, "", data, ops...)
}

// ProducerDelay  Delay延迟消息生产 delay 毫秒  5000 就是5秒
func (c *RabmqConnPool) ProducerDelay(ctx context.Context, exchangeName, routingKey string, data []byte, delay int, ops ...ProducerOption) error {
	ops = append(ops, WithProducerDelay(delay))
	return c.producerPublish(ctx, exchangeName, routingKey, data, ops...)
}

// producerPublish 生产者发送消息
func (c *RabmqConnPool) producerPublish(ctx context.Context, exchangeName, routingKey string, data []byte, ops ...ProducerOption) error {
	args := defaultProducerArgs
	for _, op := range ops {
		op(&args)
	}
	args.exchangeName = exchangeName
	args.routingKey = routingKey
	connWrapper, err := c.getConnWrapper(ctx)
	if err != nil {
		return err
	}
	channel, err := connWrapper.conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	// 确认模式
	var confirms chan amqp.Confirmation
	if args.isConfirm {
		if err := channel.Confirm(false); err != nil {
			errorx.ErrorPush("Channel could not be put into confirm mode：" + err.Error())
			return nil
		}
		confirms = channel.NotifyPublish(make(chan amqp.Confirmation, 1))
	}
	// 发送消息
	publishing := amqp.Publishing{
		ContentType:  "text/plain",
		Body:         data,
		DeliveryMode: args.publishDeliveryMode,
		Priority:     args.publishPriority,
		Headers:      amqp.Table{}, // **此处初始化 Headers**
	}

	// 处理延迟消息 // **新增代码**
	if args.delay > 0 {
		if publishing.Headers == nil {
			publishing.Headers = amqp.Table{}
		}
		publishing.Headers["x-delay"] = args.delay // **设置延迟时间**
	}

	err = channel.Publish(args.exchangeName, args.routingKey, args.publishMandatory, args.publishImmediate, publishing)
	if err != nil {
		return err
	}
	// 检查确认
	if confirms != nil {
		confirm := <-confirms
		if !confirm.Ack {
			return errors.New("消息发送失败")
		}
	}
	return nil
}
