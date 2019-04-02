package main

import (
	"flag"
	"fmt"
	"github.com/Ninlgde/rpc_go/v3.0"
)

func main() {
	addr := flag.String("addr", "127.0.0.1", "server addr")
	port := flag.String("port", "8080", "server port")

	flag.Parse()

	address := *addr + ":" + *port
	fmt.Println(address)

	v3_0.NewServer(address)
}
