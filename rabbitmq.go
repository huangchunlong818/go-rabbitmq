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

func GetNewRabbitmq(configs config.RabbitmqConfig) *RabbitmqInit {
	if newRabbitmq == nil {
		newRabbitmq = &RabbitmqInit{
			config:  &configs,
			consume: consume.NewMqConsumeServer(configs.Trace),
		}
	}
	return newRabbitmq
}

// 初始化rabbitmq 连接池，以及全局变量
func (r *RabbitmqInit) InitRabbitMq() error {
	//初始化RabbitMQX
	RabbitMQX, err := src.NewRabmqConnPool(r.config.Source)
	if err != nil {
		return errors.New("初始化rabbitmq失败：" + err.Error())
	}
	global.RabbitMq = RabbitMQX
	return nil
}

// 启动
func (r *RabbitmqInit) Start(ctx context.Context) error {
	rbConfig := r.GetRabbitmqConfig()
	if len(rbConfig) < 1 {
		return errors.New("没有定义任何rabbitmq 队列")
	}
	return r.consume.Start(ctx, rbConfig)
}

// 停止
func (r *RabbitmqInit) Stop(ctx context.Context) {
	r.consume.Stop(ctx)
}

// Direct消息生产
func (r *RabbitmqInit) Direct(ctx context.Context, key string, data []byte) error {
	now, err := r.GetRabbitmqConfigByKey(key)
	if err != nil {
		return err
	}
	return global.RabbitMq.ProducerDirect(ctx, now.Exchange, now.QueueRouterKey, data)
}

// 获取rabbitmq消费配置
func (r *RabbitmqInit) GetRabbitmqConfig(keys ...string) []config.RabbitmqConsume {
	configs := r.config.Rabbitmq

	//获取具体请求参数
	var key string
	if len(keys) > 0 {
		key = keys[0]
	}

	var rs []config.RabbitmqConsume
	if key == "" {
		//返回所有的
		for _, value := range configs {
			rs = append(rs, value)
		}
	} else {
		//返回指定的
		info, ok := configs[key]
		if ok {
			rs = append(rs, info)
		}
	}

	return rs
}

// 根据指定的 GetRabbitmqConfig 方法里的 map key 获取指定的配置
func (r *RabbitmqInit) GetRabbitmqConfigByKey(key string) (config.RabbitmqConsume, error) {
	rs := config.RabbitmqConsume{}
	data := r.GetRabbitmqConfig(key)
	if len(data) == 0 {
		return rs, errors.New("没有定义跟" + key + "有关的rabbitmq 队列")
	}
	rs = data[0]
	return rs, nil
}
