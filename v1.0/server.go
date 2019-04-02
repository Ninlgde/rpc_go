package v1_0

import (
	"net"
	"fmt"
	"encoding/json"
	"strconv"
)

type Server struct {
	Listener net.Listener
	Handlers map[string]func(conn net.Conn, params string)
}

func handler_conn(conn net.Conn, addr net.Addr, handlers map[string]func(conn net.Conn, params string))  {
	fmt.Println(addr.String(), "comes")
	for {
		length_prefix := make([]byte, 2)
		n, _ := conn.Read(length_prefix)
		if n == 0 {
			fmt.Println(addr.String(), "bye")
			conn.Close()
			return
		}
		length, _ := strconv.Atoi(string(length_prefix[:]))
		body := make([]byte, length)
		conn.Read(body)
		response := map[string]string{}
		err2 := json.Unmarshal(body, &response)
		if err2 != nil {
			fmt.Printf("Fail to unmarshal, %s\n", err2)
			return
		}
		in_ := response["in"]
		params := response["params"]
		fmt.Println(in_, params)
		handler := handlers[in_]
		handler(conn, params)
	}
}

func (s *Server)Loop() {
	for {
		conn, _ := s.Listener.Accept()
		handler_conn(conn, conn.RemoteAddr(), s.Handlers)
	}
}

func ping(conn net.Conn, params string) {
	sendresult(conn, "pong", params)
}

func sendresult(conn net.Conn, out string, params string) {
	m := map[string]string{"out":out, "params": params}
	request, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("Fail to marshal, %s\n", err)
		return
	}
	length_prefix := make([]byte, 2)
	length_prefix = []byte(strconv.Itoa(len(request)))
	conn.Write(length_prefix)
	conn.Write(request)
}