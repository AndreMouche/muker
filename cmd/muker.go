package main

import (
	"github.com/AndreMouche/logging"
	"github.com/openinx/muker/server"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	//show debug log: http://localhost:6060/debug/pprof/
	go func() {
		err := http.ListenAndServe(":6060", nil)
		if err != nil {
			logging.Warning(err)
		} else {
			logging.Info("show pprof info :http://localhost:6060/debug/pprof/")
		}

	}()

	proxySrv := server.DefaultProxyServer()

	if proxySrv == nil {
		return
	}
	// Close Proxy Backend Connection
	defer func() {
		proxySrv.Close()
	}()

	proxySrv.Start()
}
