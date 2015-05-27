package pools

import (
	"fmt"
	"github.com/openinx/muker/mysql"
)

type ConnPool struct {
	conns   chan *mysql.Client
	maxConn int
}

func NewConnPool(maxConn int) (*ConnPool, error) {
	cp := &ConnPool{
		maxConn: maxConn,
		conns:   make(chan *mysql.Client, maxConn),
	}
	for i := 0; i < maxConn; i++ {
		cli, err := mysql.NewClient()
		// Close previous opened connections if error occured
		if err != nil {
			for k := 0; k < i; k++ {
				mycli, ok := <-cp.conns
				if ok {
					mycli.Close()
				}
			}
			return cp, err
		}
		// put new connection into connection pool
		cp.conns <- cli
	}
	return cp, nil
}

func (cp *ConnPool) Get() (*mysql.Client, error) {
	select {
	case mycli, ok := <-cp.conns:
		if ok == true {
			return mycli, nil
		}
		return nil, fmt.Errorf("get conn from pool error.")
	default:
		return nil, fmt.Errorf("conn from pool has been exhaust")
	}
}

func (cp *ConnPool) Put(c *mysql.Client) error {
	select {
	case cp.conns <- c:
	default:
		return fmt.Errorf("Put conn to pool error.")
	}
	return nil
}

func (cp *ConnPool) Close() error {
	var err error
	err = nil
	for i := 0; i < len(cp.conns); i++ {
		cli, ok := <-cp.conns
		if ok {
			// Fetch the reason why close connection error.
			if err2 := cli.Close(); err2 != nil {
				err = err2
			}
		}
	}
	return err
}

func (cp *ConnPool) Cap() int {
	return cp.maxConn - len(cp.conns)
}

func (cp *ConnPool) MaxConn() int {
	return cp.maxConn
}
