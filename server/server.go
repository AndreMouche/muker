package server

import (
	"fmt"
	"net"
	"sync"
)

type ProxyServer struct {
	mu     *sync.Mutex
	Port   int
	Host   string
	ConnId uint32
}

func DefaultProxyServer() *ProxyServer {
	return &ProxyServer{
		mu:     new(sync.Mutex),
		Port:   4567,
		Host:   "127.0.0.1",
		ConnId: 0,
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
		session := NewSession(c, p.NextConnId(), 0)
		go session.HandleClient()
	}
}

func (p *ProxyServer) NextConnId() uint32 {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.ConnId++
	return p.ConnId
}
