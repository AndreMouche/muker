package server

import (
	"fmt"
	"github.com/openinx/muker/proto"
	"io"
	"net"
)

func Test() {
	fmt.Printf("Test Server Package\n")
}

func Start() {
	ln, err := net.Listen("tcp", "127.0.0.1:4567")
	if err != nil {
		fmt.Printf("Listen 127.0.0.1:4567 failed")
	}

	fmt.Printf("Listen 127.0.0.1:4567 ...\n")

	defer ln.Close()

	for {
		c, err := ln.Accept()
		if err != nil {
			fmt.Printf("Accept: %v", err)
		}

		go handleClient(c)
	}
}

func handleClient(c net.Conn) {
	defer c.Close()
	fmt.Printf("RemoteAddr: %s\n", c.RemoteAddr().String())
	var buf [1000000]byte
	n, err := c.Read(buf[:])
	if n != 0 || err != io.EOF {
		fmt.Printf("server Read = %d, %v; want 0, io.EOF\n", n, err)
		fmt.Printf("Read String: %s", string(buf[:n]))
		c.Write(proto.HandShake())
		return
	}
	c.Write([]byte("response\n"))
}
