package mysql

import (
	"fmt"
	"github.com/openinx/muker/proto"
	"github.com/openinx/muker/utils"
	"net"
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

func (c *Client) WriteCommandPacket(cmd *proto.Packet, conn net.Conn) error {
	err := c.mc.writeCommandPacketStr(cmd.Buf[0], string(cmd.Buf[1:]))
	if err != nil {
		return err
	}
	var data []byte
	var n int
	var err2 error

	// Column Length
	data, err = c.mc.readFullPacket()

	if err != nil {
		return err
	}

	fmt.Printf("read data from backend: %x\n", data)

	// Handle OK Packet, Error Packet, LoadInFile Packet.
	switch data[4] {

	// send OK Packet to client
	case iOK:
		fmt.Printf("write data to front[OK Packet]: %x\n", data)
		n, err2 = conn.Write(data)
		if err2 != nil {
			return err2
		}
		if n != len(data) {
			return fmt.Errorf("send data len is not equals.")
		}
		return nil

	//send Err Packet to client
	case iERR:
		fmt.Printf("write data to front[ERR Packet]: %x\n", data)
		n, err2 = conn.Write(data)
		if err2 != nil {
			return err2
		}
		if n != len(data) {
			return fmt.Errorf("send data len is not equals.")
		}
		return nil

	case iLocalInFile:
		return fmt.Errorf("Not support LoadInFile yet.")
	}

	//following are Column Packet.
	n, err2 = conn.Write(data)
	fmt.Printf("write data to front[Column Length Packet]: %x\n", data)
	if err2 != nil {
		return err2
	}

	columns, err3 := utils.LenEncodeToInt(data[4:])
	if err3 != nil {
		return err3
	}

	// Column Definition
	for i := uint64(0); i < columns; i++ {
		data, err = c.mc.readFullPacket()
		n, err2 = conn.Write(data)
		fmt.Printf("write data to front[Column Definition Packet]: %x\n", data)
		if err2 != nil {
			return err2
		}
		if n != len(data) {
			return fmt.Errorf("Write data error.")
		}
		// EOF Packet
		if data[4] == iEOF && (n == 9 || n == 5) {
			return nil
		}
	}

	// Column Definition EOF
	data, err = c.mc.readFullPacket()

	n, err2 = conn.Write(data)
	fmt.Printf("write data to front[Column Definition EOF]: %x\n", data)
	if err2 != nil {
		return err2
	}

	// Column Rows
	for {
		data, err = c.mc.readFullPacket()
		if err != nil {
			return err
		}
		n, err = conn.Write(data)
		fmt.Printf("write data to front: %x\n", data)
		if err != nil {
			return err
		}

		if n != len(data) {
			return fmt.Errorf("Write data error.")
		}
		// EOF Packet
		if data[4] == iEOF && (n == 9 || n == 5) {
			return nil
		}
	}

	return nil
}

func (c *Client) Close() error {
	return c.mc.Close()
}
