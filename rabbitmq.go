package rabbitmq

import (
	"context"
	"errors"
	"github.com/huangchunlong818/go-rabbitmq/rabbitmq/config"
	"github.com/huangchunlong818/go-rabbitmq/rabbitmq/consume"
	"github.com/huangchunlong818/go-rabbitmq/rabbitmq/global"
	"github.com/huangchunlong818/go-rabbitmq/rabbitmq/src"
)

type RabbitmqInit struct {
	config  *config.RabbitmqConfig
	consume *consume.MqConsumeServer
}

var newRabbitmq *RabbitmqInit

func GetNewRabbitmq(options ...config.Option) *RabbitmqInit {
	if newRabbitmq == nil {
		newRabbitmq = &RabbitmqInit{
			config: config.GetNewConfig(options...),
		}
		newRabbitmq.consume = consume.NewMqConsumeServer(newRabbitmq.config.Trace)
	}
	return newRabbitmq
}

// 初始化rabbitmq 连接池，以及全局变量
func (r *RabbitmqInit) InitRabbitMq() (*src.RabmqConnPool, error) {
	//初始化RabbitMQX
	RabbitMQX, err := src.NewRabmqConnPool(r.config.Source)
	if err != nil {
		return nil, errors.New("初始化rabbitmq失败：" + err.Error())
	}
	global.RabbitMq = RabbitMQX
	return RabbitMQX, nil
}

// 启动
func (r *RabbitmqInit) Strat(ctx context.Context) {
	r.consume.Start(ctx, r.config.Rabbitmq)
}

// 停止
func (r *RabbitmqInit) Stop(ctx context.Context) {
	r.consume.Stop(ctx)
}
