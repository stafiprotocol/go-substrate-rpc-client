package websocket_pool

import (
	"errors"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/go-substrate-rpc-client/pkg/recws"
)

var (
	ErrClosed = errors.New("pool is closed")
)

type Pool interface {
	Get() (*PoolConn, error)
	Put(*PoolConn) error
	Len() int
}

type WsPool struct {
	mu      sync.RWMutex
	exist   map[*PoolConn]bool
	conns   chan *PoolConn
	factory Factory
}

type Factory func() (*PoolConn, error)

func NewWsPool(initialCap, maxCap int, factory Factory) (Pool, error) {
	if initialCap < 0 || maxCap <= 0 || initialCap > maxCap {
		return nil, errors.New("invalid capacity settings")
	}
	c := &WsPool{
		conns:   make(chan *PoolConn, maxCap),
		exist:   map[*PoolConn]bool{},
		factory: factory,
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	for i := 0; i < initialCap; i++ {
		conn, err := factory()
		if err != nil {
			c.Close()
			return nil, fmt.Errorf("factory is not able to fill the pool: %s", err)
		}
		c.conns <- conn
		c.exist[conn] = true
	}
	return c, nil
}

func (c *WsPool) getConnsAndFactory() (chan *PoolConn, Factory) {
	conns := c.conns
	factory := c.factory
	return conns, factory
}

func (c *WsPool) Get() (*PoolConn, error) {
	conns, factory := c.getConnsAndFactory()
	if conns == nil {
		return nil, ErrClosed
	}
	logrus.Tracef("WsPool.Get pool len :%d", len(c.conns))

	var err error
	select {
	case conn := <-conns:
		if conn == nil || !conn.Conn.IsConnected() {
			c.mu.Lock()
			delete(c.exist, conn)
			c.mu.Unlock()

			logrus.Trace("use factory reconnect")
			conn, err = factory()
			if err != nil {
				return nil, err
			}
		} else {
			c.mu.Lock()
			delete(c.exist, conn)
			c.mu.Unlock()
		}
		return conn, nil
	default:
		conn, err := factory()
		if err != nil {
			return nil, err
		}
		return conn, nil
	}
}

func (c *WsPool) Put(conn *PoolConn) error {
	if conn == nil {
		return errors.New("connection is nil. rejecting")
	}
	conn.mu.Lock()
	defer conn.mu.Unlock()

	logrus.Tracef("WsPool.Put pool len :%d", len(c.conns))
	if conn.unusable {
		if conn.Conn != nil {
			conn.Conn.Close()
		}
		return nil
	}

	c.mu.RLock()
	if c.exist[conn] {
		c.mu.RUnlock()
		return nil
	}
	c.mu.RUnlock()

	if c.conns == nil {
		conn.Conn.Close()
		return nil
	}

	select {
	case c.conns <- conn:
		c.mu.Lock()
		c.exist[conn] = true
		c.mu.Unlock()
		return nil
	default:
		conn.Conn.Close()
		return nil
	}
}

func (c *WsPool) Close() {
	c.mu.Lock()
	conns := c.conns
	c.conns = nil
	c.factory = nil
	c.mu.Unlock()

	if conns == nil {
		return
	}

	close(conns)
	for conn := range conns {
		conn.Conn.Close()
	}
}

func (c *WsPool) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	conns, _ := c.getConnsAndFactory()
	return len(conns)
}

type PoolConn struct {
	Conn     *recws.RecConn
	mu       sync.RWMutex
	unusable bool
}

// MarkUnusable() marks the connection not usable any more, to let the pool close it instead of returning it to pool.
func (p *PoolConn) MarkUnusable() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.unusable = true
}

// newConn wraps a standard net.Conn to a poolConn net.Conn.
func WrapConn(conn *recws.RecConn) *PoolConn {
	p := &PoolConn{}
	p.Conn = conn
	return p
}
