package src

import (
	"context"
	"errors"
	errorx "github.com/huangchunlong818/go-rabbitmq/rabbitmq/error"
	"github.com/rs/xid"
	"github.com/streadway/amqp"
	"sync"
	"sync/atomic"
	"time"
)

const (
	STATUS_CONNECTING      = 1               // 表示连接中
	STATUS_DISCONNECTED    = 2               // 表示连接断开
	DEFAULT_MAX_IDLECONNS  = 2               // 默认最大连接数
	RETRY_GET_CONNDURATION = 5 * time.Second // 默认获取连接失败后重试间隔
)

// mqConnWrapper mq连接和使用状态
type mqConnWrapper struct {
	conn        *amqp.Connection // AMQP连接
	status      int32            // 连接状态; 1: 连接中; 2：连接断开
	useCount    int32            // 使用计数
	index       string           // 连接池中的索引
	closedIndex chan string
}

// notifyClose 监听连接断开事件
func (c *mqConnWrapper) notifyClose() {
	receiver := make(chan *amqp.Error)
	notifyClose := c.conn.NotifyClose(receiver)
	amqpErr := <-notifyClose
	if amqpErr != nil {
		errorx.ErrorPush("AMQP connection closed：" + amqpErr.Error())
	}
	atomic.StoreInt32(&c.status, STATUS_DISCONNECTED)
	c.conn.Close()
	c.closedIndex <- c.index
}

// checkClosed 检查连接是否断开
func (c *mqConnWrapper) checkClosed() bool {
	return atomic.LoadInt32(&c.status) == STATUS_DISCONNECTED
}

// useAdd 使用计数+1
func (c *mqConnWrapper) useAdd() {
	atomic.AddInt32(&c.useCount, 1)
}

// 获取使用计数
func (c *mqConnWrapper) getUseCount() int32 {
	return atomic.LoadInt32(&c.useCount)
}

// RabmqConnPool 是一个RabbitMQ连接池结构体
type RabmqConnPool struct {
	source       string           // RabbitMQ地址
	mu           sync.Mutex       // 互斥锁
	conns        []*mqConnWrapper // 连接池
	maxIdleConns int              // 连接池最大空闲连接数
	closedIndex  chan string      // 已关闭的连接索引
}

// PoolOptinon 连接池配置选项
type PoolOptinon func(*RabmqConnPool)

// WithIdleTimeout 设置连接池中最大连接数
func WithMaxIdleConns(maxIdleConns int) PoolOptinon {
	return func(pool *RabmqConnPool) {
		pool.maxIdleConns = maxIdleConns
	}
}

// NewRabmqConnPool 创建一个连接池
func NewRabmqConnPool(source string, connOpts ...PoolOptinon) (*RabmqConnPool, error) {
	pool := &RabmqConnPool{
		source:       source,
		maxIdleConns: DEFAULT_MAX_IDLECONNS,
		conns:        make([]*mqConnWrapper, 0),
		closedIndex:  make(chan string, 10),
	}
	// 设置连接池配置选项
	for _, o := range connOpts {
		o(pool)
	}
	// 初始化连接池
	for i := 0; i < pool.maxIdleConns; i++ {
		connWrap, err := pool.newConnWrapper()
		if err != nil {
			return nil, err
		}
		pool.conns = append(pool.conns, connWrap)
	}
	// 启动连接池监控
	go pool.monitor()
	return pool, nil
}

// newConnWrapper 创建连接wrapper
func (c *RabmqConnPool) newConnWrapper() (*mqConnWrapper, error) {
	conn, err := amqp.Dial(c.source)
	if err != nil {
		return nil, err
	}
	connWrap := &mqConnWrapper{
		conn:        conn,
		status:      STATUS_CONNECTING,
		index:       c.GenXID(),
		closedIndex: c.closedIndex,
	}
	go connWrap.notifyClose()
	return connWrap, nil
}

// GenXID 生成XID
func (c *RabmqConnPool) GenXID() string {
	return xid.New().String()
}

// getConnWrapper 从连接池中获取AMQP连接，选择引用数最小的连接
func (c *RabmqConnPool) getConnWrapper(ctx context.Context) (*mqConnWrapper, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var res *mqConnWrapper
	for _, cw := range c.conns {
		if cw == nil {
			continue // 跳过nil值
		}
		if cw.checkClosed() {
			continue
		}
		if res == nil {
			res = cw
			continue
		}
		if cw.getUseCount() < res.getUseCount() {
			res = cw
		}
	}
	if res != nil {
		res.useAdd()
		return res, nil
	}
	// 如果连接池中没有可用的连接，则创建新的连接
	if len(c.conns) >= c.maxIdleConns {
		return nil, errors.New("连接池中没有可用的连接")
	}
	connWrap, err := c.newConnWrapper()
	if err != nil {
		return nil, err
	}
	connWrap.useAdd()
	c.conns = append(c.conns, connWrap)
	return connWrap, nil
}

// monitor 监控连接池中的连接，并删除已断开的连接
func (c *RabmqConnPool) monitor() {
	for index := range c.closedIndex {
		c.mu.Lock()
		// 删除已断开的连接
		for i, cw := range c.conns {
			if cw.index == index {
				c.conns = append(c.conns[:i], c.conns[i+1:]...)
				break
			}
		}
		// 如果连接池中的连接数小于最大连接数，则创建新的连接
		if len(c.conns) >= c.maxIdleConns {
			continue
		}
		connWrap, err := c.newConnWrapper()
		if err != nil {
			errorx.ErrorPush("创建新的AMQP连接失败：" + err.Error())
		}
		c.conns = append(c.conns, connWrap)
		c.mu.Unlock()
	}
}
