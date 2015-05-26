package mysql

import (
	"fmt"
	"github.com/openinx/muker/proto"
	"io"
)

type Client struct {
	mc *mysqlConn
}

func NewClient() (*Client, error) {
	md := &MySQLDriver{}
	mc, err := md.Open("root:123456@tcp(127.0.0.1:3306)/test?autocommit=true")
	if err != nil {
		return nil, err
	}

	c := &Client{mc: mc}
	return c, nil
}

func (c *Client) WriteCommandPacket(cmd *proto.Packet, w io.Writer) error {
	err := c.mc.writeCommandPacketStr(cmd.Buf[0], string(cmd.Buf[1:]))
	if err != nil {
		return err
	}
	var data []byte

	for {
		data, err = c.mc.readPacket()
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			return err
		}
		fmt.Printf("recv data <-- backend: %x\n", data)
		n, err2 := w.Write(data)
		if err2 != nil {
			return err2
		}
		if n != len(data) {
			return fmt.Errorf("Write data error.")
		}
	}
}
