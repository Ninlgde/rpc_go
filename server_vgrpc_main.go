package main

import (
	"flag"
	"fmt"
	"github.com/Ninlgde/rpc_go/vgrpc"
)

func main() {
	addr := flag.String("addr", "127.0.0.1", "server addr")
	port := flag.String("port", "8080", "server port")

	flag.Parse()

	address := *addr + ":" + *port
	fmt.Println(address)

	vgrpc.NewServer(address)
}
