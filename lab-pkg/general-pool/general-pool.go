package main

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
	"ultron-lab/lab-pkg/backoff"

	log "github.com/sirupsen/logrus"
)

var (
	ErrMaxActiveConnReached = errors.New("MaxActiveConnReached")
	ErrClosed               = errors.New("pool is closed")
)

// Pool 基本方法
type Pool interface {
	// 获取连接
	Get() (interface{}, error)
	// 放回连接
	Put(interface{}) error
	// 关闭连接
	Close(interface{}) error
	// 释放所有连接
	Release()
	// 获取连接池数量
	Len() int
}

// Config 连接池相关配置
type Config struct {
	//连接池中拥有的最小连接数
	InitialCap int
	//最大并发存活连接数
	MaxCap int
	//最大空闲连接
	MaxIdle int
	//生成连接的方法
	Factory func() (interface{}, error)
	//关闭连接的方法
	Close func(interface{}) error
	//检查连接是否有效的方法
	Ping func(interface{}) error
	//连接最大空闲时间，超过该事件则将失效
	IdleTimeout time.Duration
}

// channelPool 存放连接信息
type channelPool struct {
	mu                       sync.RWMutex
	conns                    chan *idleConn
	factory                  func() (interface{}, error)
	close                    func(interface{}) error
	ping                     func(interface{}) error
	idleTimeout, waitTimeOut time.Duration
	maxActive                int
	openingConns             int
}

type idleConn struct {
	conn interface{}
	t    time.Time
}

// NewChannelPool 初始化连接
func NewChannelPool(poolConfig *Config) (Pool, error) {
	if !(poolConfig.InitialCap <= poolConfig.MaxIdle && poolConfig.MaxCap >= poolConfig.MaxIdle && poolConfig.InitialCap >= 0) {
		return nil, errors.New("invalid capacity settings")
	}
	if poolConfig.Factory == nil {
		return nil, errors.New("invalid factory func settings")
	}
	if poolConfig.Close == nil {
		return nil, errors.New("invalid close func settings")
	}

	c := &channelPool{
		conns:        make(chan *idleConn, poolConfig.MaxIdle),
		factory:      poolConfig.Factory,
		close:        poolConfig.Close,
		idleTimeout:  poolConfig.IdleTimeout,
		maxActive:    poolConfig.MaxCap,
		openingConns: poolConfig.InitialCap,
	}

	if poolConfig.Ping != nil {
		c.ping = poolConfig.Ping
	}

	for i := 0; i < poolConfig.InitialCap; i++ {
		conn, err := c.factory()
		time.Sleep(1 * time.Second)
		if err != nil {
			c.Release()
			return nil, fmt.Errorf("factory is not able to fill the pool: %s", err)
		}
		c.conns <- &idleConn{conn: conn, t: time.Now()}
	}

	return c, nil
}

// getConns 获取所有连接
func (c *channelPool) getConns() chan *idleConn {
	c.mu.Lock()
	conns := c.conns
	c.mu.Unlock()
	return conns
}

// Get 从pool中取一个连接
func (c *channelPool) Get() (interface{}, error) {
	conns := c.getConns()
	if conns == nil {
		return nil, ErrClosed
	}
	for {
		select {
		case wrapConn := <-conns:
			if wrapConn == nil {
				return nil, ErrClosed
			}
			//判断是否超时，超时则丢弃
			if timeout := c.idleTimeout; timeout > 0 {
				if wrapConn.t.Add(timeout).Before(time.Now()) {
					//丢弃并关闭该连接
					c.Close(wrapConn.conn)
					continue
				}
			}
			//判断是否失效，失效则丢弃，如果用户没有设定 ping 方法，就不检查
			if c.ping != nil {
				if err := c.Ping(wrapConn.conn); err != nil {
					c.Close(wrapConn.conn)
					continue
				}
			}
			return wrapConn.conn, nil
		default:
			c.mu.Lock()
			log.Debugf("openConn %v %v", c.openingConns, c.maxActive)
			defer c.mu.Unlock()
			if c.openingConns >= c.maxActive {
				return nil, ErrMaxActiveConnReached
			}
			if c.factory == nil {
				return nil, ErrClosed
			}
			conn, err := c.factory()
			if err != nil {
				return nil, err
			}
			c.openingConns++
			return conn, nil
		}
	}
}

// Put 将连接放回pool中
func (c *channelPool) Put(conn interface{}) error {
	if conn == nil {
		return errors.New("connection is nil. rejecting")
	}

	c.mu.Lock()

	if c.conns == nil {
		c.mu.Unlock()
		return c.Close(conn)
	}

	select {
	case c.conns <- &idleConn{conn: conn, t: time.Now()}:
		c.mu.Unlock()
		return nil
	default:
		c.mu.Unlock()
		//连接池已满，直接关闭该连接
		return c.Close(conn)
	}
}

// Close 关闭单条连接
func (c *channelPool) Close(conn interface{}) error {
	if conn == nil {
		return errors.New("connection is nil. rejecting")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.close == nil {
		return nil
	}
	c.openingConns--
	return c.close(conn)
}

// Ping 检查单条连接是否有效
func (c *channelPool) Ping(conn interface{}) error {
	if conn == nil {
		return errors.New("connection is nil. rejecting")
	}
	return c.ping(conn)
}

// Release 释放连接池中所有连接
func (c *channelPool) Release() {
	c.mu.Lock()
	conns := c.conns
	c.conns = nil
	c.factory = nil
	c.ping = nil
	closeFun := c.close
	c.close = nil
	c.mu.Unlock()

	if conns == nil {
		return
	}

	close(conns)
	for wrapConn := range conns {
		//log.Printf("Type %v\n",reflect.TypeOf(wrapConn.conn))
		closeFun(wrapConn.conn)
	}
}

// Len 连接池中已有的连接
func (c *channelPool) Len() int {
	return len(c.getConns())
}

func main() {
	// init backoff
	b := backoff.NewBackOff(
		backoff.WithJitterFlag(true),
		backoff.WithMaxDelay(120*time.Second))
	done := make(chan struct{}, 1)
	poolConfig := Config{
		InitialCap: 1,
		MaxCap:     1,
		MaxIdle:    1,
		Factory: func() (conn interface{}, err error) {
			conn, err = net.Dial("tcp", "10.40.9.80:38888")
			log.Infof("建立连接, @%v", time.Now().String())
			return
		},
		Close: func(conn interface{}) error {
			return conn.(net.Conn).Close()
		},
		Ping: func(conn interface{}) error {
			_, err := conn.(net.Conn).Write([]byte("哈哈哈\n"))
			return err
		},
		IdleTimeout: 30 * time.Second,
	}

	pool, err := NewChannelPool(&poolConfig)

	if err != nil {
		log.Error(err)
		return
	}

	go func(pool Pool) {
		for i := 0; i < 1000; i++ {
			time.Sleep(1 * time.Second)
			log.Infof("pool size: %d", pool.Len())
		}
		done <- struct{}{}
	}(pool)

	go func(pool Pool) {
		for i := 0; i < 100; i++ {
			time.Sleep(1 * time.Second)
			conn, err := pool.Get()
			if err != nil {
				log.Error(err)
				// wait
				b.Sleep()
				continue
			}
			strToSend := fmt.Sprintf("time: %s, data: fuck\n", time.Now().String())
			_, err = conn.(net.Conn).Write([]byte(strToSend))

			if err != nil {
				log.Error(err)
			}

			pool.Put(conn)
		}
	}(pool)

	<-done
}
