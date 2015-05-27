package main

import (
	"github.com/openinx/muker/server"
)

func main() {
	proxySrv := server.DefaultProxyServer()

	// Close Proxy Backend Connection
	defer func() {
		proxySrv.Close()
	}()

	proxySrv.Start()
}
