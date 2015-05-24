package server

import (
	"net"
)

type DBConn struct {
	conn         net.Conn
	connectionId uint32
	sequenceId   uint8
}

type FrontDBConn struct {
}

type BackendConn struct {
}
