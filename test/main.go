package main

import (
	"context"
	"fmt"
	"github.com/huangchunlong818/go-rabbitmq"
	"github.com/huangchunlong818/go-rabbitmq/rabbitmq/config"
)

func ttt(ctx context.Context, data []byte) error {
	fmt.Println(data)
	return nil
}

func main() {
	//配置
	ctx := context.Background()
	// Rabbitmq 队列配置
	var Configs = config.RabbitmqConfig{
		Rabbitmq: map[string]config.RabbitmqConsume{
			//test
			"common": {
				Operation:      "/rabbitmq/tt/ttttttttttt", //定义rabbitMQ消费自定义路径
				Exchange:       "test-delay",               //定义消费rabbitmq的交换器
				Queue:          "test-delay-queue",         //定义rabbitmq的队列名称
				QueueRouterKey: "test-delay-queue",         //定义rabbitmq的路由键
				Handler:        ttt,                        //定义真正的消费逻辑，必须是 func (ctx context.Context, data []byte) error 类型
				ConsumerNum:    2,                          //定义几个协程消费，如果消费次数大于1，则最好定义2以上，因为重试的时候阻塞等待间隔时间
			},
		},
		Trace: config.TraceConfig{
			Enable: false, //开启trace
		},
		Source: "amqp://huangchunlong:5R84Y8vtzyfzqTmbQ@192.168.15.40:5672/", //rabbitmq
	}

	// 获取操作对象
	mq := rabbitmq.GetNewRabbitmq(Configs)

	//初始化连接池
	err := mq.InitRabbitMq()
	if err != nil {
		fmt.Println(err)
	}

	//操作rabbitmq-写入, key 对应的上面的配置中的 Rabbitmq 下面的配置的key
	err = mq.Delay(ctx, "common", []byte("testValue"), 20000)
	if err != nil {
		fmt.Println(err)
	}

	////消费服务
	//mq.Start(ctx)
	//
	////停止消费服务
	//mq.Stop(ctx)
	//
	//fmt.Println("Hello World", ctx, mq)
}
