package main

import (
	"context"
	"fmt"
	"github.com/huangchunlong818/go-rabbitmq"
	"github.com/huangchunlong818/go-rabbitmq/rabbitmq/config"
)

func main() {
	//配置
	ctx := context.Background()
	var rbConfig []config.RabbitmqConsume
	rbConfig = append(rbConfig, config.RabbitmqConsume{
		Operation:      "",
		Exchange:       "",
		Queue:          "",
		QueueRouterKey: "",
		Handler:        nil,
		ConsumerNum:    0,
		Num:            0,
		TimeInterval:   0,
	})
	traceC := config.TraceConfig{
		Enable: false, //不开启trace
	}

	// 获取操作对象
	mq := rabbitmq.GetNewRabbitmq(config.WithSource("amqp://connector:connector123@192.168.15.40:5672/connector"), config.WithRabbitmq(rbConfig), config.WithTrace(traceC))

	//初始化连接池
	Rabbitmq, err := mq.InitRabbitMq()
	if err != nil {
		fmt.Println(err)
	}

	//操作rabbitmq-写入
	err = Rabbitmq.ProducerDirect(ctx, "testChange", "testKey", []byte("testValue"))
	if err != nil {
		fmt.Println(err)
	}

	//消费服务
	mq.Strat(ctx)

	fmt.Println("Hello World", ctx, mq)
}
