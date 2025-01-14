package config

import (
	"context"
	"time"
)

type RabbitmqConfig struct {
	Rabbitmq map[string]RabbitmqConsume
	Trace    TraceConfig
	Source   string
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
	Type         string        //类型， delay 延迟队列，direct 直接队列普通队列,目前只支持这2种，如果没定义，默认都是 direct
}

// consumefunc 数据处理
type Consumefunc func(ctx context.Context, data []byte) error
