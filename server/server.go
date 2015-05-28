package server

import (
	"fmt"
	"github.com/openinx/muker/pools"
	"net"
	"sync"
)

type ProxyServer struct {
	mu          *sync.Mutex
	Port        int
	Host        string
	ConnId      uint32
	BackendPool *pools.ConnPool
}

func DefaultProxyServer() *ProxyServer {

	// Initialize 50 conn pools
	backends, err := pools.NewConnPool(50)
	if err != nil {
		fmt.Printf("open backend connection pool failed: %v", err)
	}

	return &ProxyServer{
		mu:          new(sync.Mutex),
		Port:        4567,
		Host:        "127.0.0.1",
		ConnId:      0,
		BackendPool: backends,
	}
}

func (p *ProxyServer) Start() {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", p.Host, p.Port))
	if err != nil {
		fmt.Printf("Listen %s:%d failed\n", p.Host, p.Port)
	}

	fmt.Printf("Listen %s:%d ...\n", p.Host, p.Port)

	defer ln.Close()

	for {
		c, err := ln.Accept()
		if err != nil {
			fmt.Printf("Accept error: %s\n", err.Error())
		}
		session := NewSession(c, p.NextConnId(), 0, p.BackendPool)
		go session.HandleClient()
	}
}

func (p *ProxyServer) NextConnId() uint32 {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.ConnId++
	return p.ConnId
}

func (p *ProxyServer) Close() {
	err := p.BackendPool.Close()
	if err != nil {
		fmt.Printf("close backend connection pool failed: %v", err)
	}
}
