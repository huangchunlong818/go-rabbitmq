package main

import (
	"context"
	"fmt"
	"github.com/huangchunlong818/go-rabbitmq"
	"github.com/huangchunlong818/go-rabbitmq/rabbitmq/config"
	"time"
)

func main() {
	//配置
	ctx := context.Background()
	// Rabbitmq 队列配置
	var Configs = config.RabbitmqConfig{
		Rabbitmq: map[string]config.RabbitmqConsume{
			//test
			"common": {
				Operation:      "/rabbitmq/shopify/commonHandler", //定义rabbitMQ消费自定义路径
				Exchange:       "connector_direct_exchange",       //定义消费rabbitmq的交换器
				Queue:          "shopify_common_router",           //定义rabbitmq的队列名称
				QueueRouterKey: "shopify_common_router",           //定义rabbitmq的路由键
				Handler:        nil,                               //定义真正的消费逻辑，必须是 func (ctx context.Context, data []byte) error 类型
				ConsumerNum:    2,                                 //定义几个协程消费，如果消费次数大于1，则最好定义2以上，因为重试的时候阻塞等待间隔时间
				Num:            2,                                 //定义消费次数
				TimeInterval:   10 * time.Second,                  //定义每次消费重试间隔时间，阻塞等待
			},
		},
		Trace: config.TraceConfig{
			Enable:           true, //开启trace
			Version:          "",   //版本号，自定义
			Endpoint:         "",   //阿里云Endpoint
			Project:          "",   //阿里云Project
			Instance:         "",   //阿里云Instance
			AccessKeyId:      "",   //阿里云AccessKeyId
			AccessKeySecret:  "",   //阿里云AccessKeySecret
			ServiceName:      "",   //当前服务定义的服务名称，自定义
			ServiceNamespace: "",   //当前服务所在服务集群域名，自定义
		},
		Source: "amqp://connector:connector123@192.168.15.40:5672/connector", //rabbitmq
	}

	// 获取操作对象
	mq := rabbitmq.GetNewRabbitmq(Configs)

	//初始化连接池
	err := mq.InitRabbitMq()
	if err != nil {
		fmt.Println(err)
	}

	//操作rabbitmq-写入, key 对应的上面的配置中的 Rabbitmq 下面的配置的key
	err = mq.Direct(ctx, "common", []byte("testValue"))
	if err != nil {
		fmt.Println(err)
	}

	//消费服务
	mq.Start(ctx)

	//停止消费服务
	mq.Stop(ctx)

	fmt.Println("Hello World", ctx, mq)
}
