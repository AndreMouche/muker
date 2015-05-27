package main

import (
	"github.com/openinx/muker/server"
)

func main() {
	proxySrv := server.DefaultProxyServer()
	proxySrv.Start()
}
