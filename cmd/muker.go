package main

import (
	"fmt"
	"github.com/openinx/muker/client"
	"github.com/openinx/muker/proto"
	"github.com/openinx/muker/server"
	"github.com/openinx/muker/utils"
)

func main() {
	x := proto.HandShake()
	fmt.Printf("len: %d\n", x)
	client.TestClient()
	utils.Hello()
	server.Start()
}
