package v1_0

import (
	"net"
	"fmt"
	"testing"
)

func TestServer(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Printf("Fail to connect, %s\n", err)
	}
	defer listener.Close()
	handlers := make(map[string]func(conn net.Conn, params string))
	handlers["ping"] = ping
	server := &Server{
		Listener:listener,
		Handlers:handlers,
	}
	server.Loop()
}