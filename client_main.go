package main

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/Ninlgde/rpc_go/v1.0"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8081")
	if err != nil {
		fmt.Printf("Fail to connect, %s\n", err)
		return
	}
	defer conn.Close()
	client := v1_0.Client{conn}
	for i := 0; i < 10; i++ {
		out, result := client.Rpc("ping", "ireader "+strconv.Itoa(i))
		fmt.Println(out, result)
		time.Sleep(time.Second * 1)
	}
}
