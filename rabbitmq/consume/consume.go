package consume

import (
	"context"
	"errors"
	"fmt"
	"github.com/huangchunlong818/go-rabbitmq/rabbitmq/config"
	errorx "github.com/huangchunlong818/go-rabbitmq/rabbitmq/error"
	"github.com/huangchunlong818/go-rabbitmq/rabbitmq/global"
	rabbitmqsrc "github.com/huangchunlong818/go-rabbitmq/rabbitmq/src"
	tracee "github.com/huangchunlong818/go-rabbitmq/rabbitmq/trace"
	traces "github.com/huangchunlong818/go-trace/trace"
	"go.opentelemetry.io/otel/trace"
	"sync/atomic"
	"time"
)

const (
	MQ_SERVER_STATUS_RUNNING = 1 //  运行状态
	MQ_SERVER_STATUS_CLOSED  = 2 // 关闭状态
)

type MqConsumeServer struct {
	isClosed    uint32             // 服务运行状态
	isTrace     uint32             // trace启用禁用
	traceConfig config.TraceConfig //trace相关配置
}

// NewMqConsumeServer 初始化MqConsumeServer
func NewMqConsumeServer(traceConfig config.TraceConfig) *MqConsumeServer {
	return &MqConsumeServer{
		isClosed:    MQ_SERVER_STATUS_RUNNING,
		traceConfig: traceConfig,
	}
}

// Start 启动服务
func (c *MqConsumeServer) Start(ctx context.Context, configs []config.RabbitmqConsume) error {

	//目前暂时只扩展了 Direct 交换器
	for _, conf := range configs {
		c.registerConsumeDirectHandler(conf.Operation, conf.Exchange, conf.Queue, conf.QueueRouterKey, conf.Handler, conf.ConsumerNum, conf.Num, conf.TimeInterval)
	}

	return nil
}

// Stop 停止服务
func (c *MqConsumeServer) Stop(ctx context.Context) error {
	atomic.StoreUint32(&c.isClosed, MQ_SERVER_STATUS_CLOSED)
	return nil
}

// checkClosed 检查是否关闭
func (c *MqConsumeServer) checkClosed() bool {
	return atomic.LoadUint32(&c.isClosed) == MQ_SERVER_STATUS_CLOSED
}

// registerConsumeDirectHandler direct消费处理
func (c *MqConsumeServer) registerConsumeDirectHandler(operation string, exchangeName, queueName, routingKey string, handler config.Consumefunc, consumerNum int, num int, timeInterval time.Duration, ops ...rabbitmqsrc.ConsumeOption) {
	for i := 0; i < consumerNum; i++ {
		//根据配置的消费协程数量，开启多少个协程 去消费
		go func() {
			defer func() {
				err := recover()
				if err != nil {
					errorx.ErrorPush("队列1：" + queueName + " ,operation:" + operation + "，错误信息：" + fmt.Errorf("%+v", err).Error())
				}
			}()
			global.RabbitMq.ConsumeDirect(exchangeName, queueName, routingKey, func(data []byte) error {
				return c.consumeHandle(operation, data, handler)
			}, num, timeInterval, ops...)
		}()
	}
}

// consumeHandle 带链路跟踪处理
func (c *MqConsumeServer) consumeHandle(operation string, data []byte, handler config.Consumefunc) error {
	if c.checkClosed() {
		return errors.New("------closed-----")
	}

	var (
		ctx       = context.Background() //空上下文
		traceSpan trace.Span             //trace组件的 span
		tr        = true                 //是否真正上传trace
		tracerObj *traces.TracerSpan     //组件封装的 trace span
	)
	if c.traceConfig.Enable {
		if tracer, err := tracee.GetTraces(c.traceConfig); err != nil {
			errorx.ErrorPush("数据：" + string(data) + " ,operation:" + operation + "，错误信息：" + fmt.Errorf("%+v", err).Error() + "，启用trace失败，跳过trace追踪")
			tr = false
		} else {
			tracerObj = tracer
			ctxs, panicSpan, errx := tracerObj.StartSpan(ctx, operation, "rabbitmq")
			if errx != nil {
				errorx.ErrorPush("数据：" + string(data) + " ,operation:" + operation + "，错误信息：" + fmt.Errorf("%+v", err).Error() + "，启用traceSpan失败，跳过trace追踪")
				tr = false
			} else {
				traceSpan = panicSpan
				ctx = ctxs
			}
		}
	}
	err := handler(ctx, data)

	defer func() {
		if tr {
			tracerObj.EndSpan(ctx, traceSpan, err)
		}
	}()

	return err
}
