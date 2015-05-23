package server

import (
	"errors"
	"fmt"
	"github.com/openinx/muker/proto"
	"github.com/openinx/muker/utils"
	"io"
	"net"
	"time"
)

func Test() {
	fmt.Printf("Test Server Package\n")
}

func Start() {
	ln, err := net.Listen("tcp", "127.0.0.1:4567")
	if err != nil {
		fmt.Printf("Listen 127.0.0.1:4567 failed\n")
	}

	fmt.Printf("Listen 127.0.0.1:4567 ...\n")

	defer ln.Close()

	for {
		c, err := ln.Accept()
		if err != nil {
			fmt.Printf("Accept error: %s\n", err.Error())
		}

		go handleClient(c)
	}
}

func handleClient(c net.Conn) {
	fmt.Printf("RemoteAddr: %s\n", c.RemoteAddr().String())
	c.Write(proto.HandShake())

	for {
		p, err := readPacket(c)
		if err == io.EOF {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		if err != nil && err != io.EOF {
			fmt.Printf(err.Error() + "\n")
			c.Close()
			return
		}
		handlePacket(c, p)
	}
}

func handlePacket(c net.Conn, p *proto.Packet) {
	fmt.Printf("send to client packet : %x\n", p.ToBytes())
	c.Write(proto.OK())
}

func readPacket(c net.Conn) (*proto.Packet, error) {
	header := make([]byte, 4)
	n, err := c.Read(header[:])

	if err != nil {
		return nil, err
	}

	if n != 4 {
		fmt.Printf("header length: %d\n", n)
		return nil, errors.New("Read packet header error : less than 4 bytes")
	}

	fmt.Printf("Read header: %x\n", header)

	packetLength := utils.BytesToUint24(header[:3])
	sequenceId := utils.BytesToUint8(header[3:])

	fmt.Printf("Read PacketLength: %d\n", packetLength)

	body := make([]byte, packetLength)
	n, err = c.Read(body)

	if err != nil && err != io.EOF {
		return nil, err
	}

	if n != packetLength {
		return nil, errors.New(fmt.Sprintf("Read packet body error: head length: %d, body length: %d", packetLength, n))
	}

	return proto.NewPacket(sequenceId, body), nil
}
