package main

import (
	"fmt"
	"github.com/openinx/muker/server"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	//show debug log: http://localhost:6060/debug/pprof/
	go func() {
		rerr := http.ListenAndServe(":6060", nil)
		if rerr != nil {
			fmt.Printf("%v\n", rerr)
		}
	}()

	proxySrv := server.DefaultProxyServer()

	// Close Proxy Backend Connection
	defer func() {
		proxySrv.Close()
	}()

	proxySrv.Start()
}
