package config

import (
	"context"
	"time"
)

type RabbitmqConfig struct {
	Rabbitmq []RabbitmqConsume
	Trace    TraceConfig
	Source   string
}

var newRabbitmqConfig *RabbitmqConfig

func GetNewConfig(options ...Option) *RabbitmqConfig {
	if newRabbitmqConfig == nil {
		newRabbitmqConfig = &RabbitmqConfig{}
	}

	for _, option := range options {
		option(newRabbitmqConfig)
	}

	return newRabbitmqConfig
}

type Option func(*RabbitmqConfig)

func WithTrace(t TraceConfig) Option {
	return func(config *RabbitmqConfig) {
		config.Trace = t
	}
}

func WithSource(t string) Option {
	return func(config *RabbitmqConfig) {
		config.Source = t
	}
}

func WithRabbitmq(t []RabbitmqConsume) Option {
	return func(config *RabbitmqConfig) {
		config.Rabbitmq = t
	}
}

// trace配置
type TraceConfig struct {
	Enable           bool
	Version          string
	Endpoint         string
	Project          string
	Instance         string
	AccessKeyId      string
	AccessKeySecret  string
	ServiceName      string
	ServiceNamespace string
}

// 定义rabbitmq 消费服务启动，需要的配置结构
type RabbitmqConsume struct {
	Operation      string      //消费路径，自定义，可用于内部追踪
	Exchange       string      //交换器
	Queue          string      //队列名称
	QueueRouterKey string      //队列键
	Handler        Consumefunc //消费队列数据处理方法
	ConsumerNum    int         //消费当前队列数量， 最好2个或以上，因为消费等待是阻塞的
	//消费多次的情况下，每次消费时间间隔，默认30秒，这是阻塞消费协程的，谨慎配置，
	Num          int           //消费次数，最低是1次哦，不配置或者配置0 ，默认1
	TimeInterval time.Duration //消费多次的情况下，每次消费时间间隔，默认30秒
}

// consumefunc 数据处理
type Consumefunc func(ctx context.Context, data []byte) error
