package v2_0

import (
	"fmt"
	"net"
	"strconv"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Printf("Fail to connect, %s\n", err)
		return
	}
	defer conn.Close()
	client := Client{conn}
	for i := 0; i < 10; i++ {
		out, result := client.Rpc("ping", "ireader "+strconv.Itoa(i))
		fmt.Println(out, result)
		time.Sleep(time.Second * 1)
	}
}
