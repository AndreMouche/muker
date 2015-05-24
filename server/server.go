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

	// send server information
	c.Write(proto.HandShake())

	// auth mysql-client
	for {
		p, err := readPacket(c)
		fmt.Printf("auth: mysql-client\n")
		if err == io.EOF {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		c.Write(proto.OK(p.SeqId + 1))
		break
	}

	// command phrase
	for {
		p, err := readPacket(c)
		if err == io.EOF {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		if err != nil {
			fmt.Printf(err.Error() + "\n")
			c.Close()
			return
		}
		handleCommand(c, p)
	}
}

func handleCommand(c net.Conn, p *proto.Packet) {
	fmt.Printf("recv client packet : %x\n", p.ToBytes())
	comType := p.Body[0]
	if comType == proto.ComSleep {
	} else if comType == proto.ComQuit {
		c.Write(proto.OK(p.SeqId + 1))
	} else if comType == proto.ComQuery {
		fmt.Printf("Query Command: %s\n", string(p.Body[1:]))
		c.Write(proto.OK(p.SeqId + 1))
	} else if comType == proto.ComPing {
		fmt.Printf("Command Ping\n")
		c.Write(proto.OK(p.SeqId + 1))
	} else if comType == proto.ComInitDB {
		schemaName := string(p.Body[1:])
		fmt.Printf("Command Init DB: %s\n", schemaName)
		c.Write(proto.OK(p.SeqId + 1))
	} else if comType == proto.ComCreateDB {
		c.Write(proto.OK(p.SeqId + 1))
	}
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
