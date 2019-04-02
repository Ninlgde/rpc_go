package v2_0

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type Server struct {
	Listener net.Listener
	Handlers map[string]func(conn net.Conn, params string)
}

func handler_conn(conn net.Conn, addr net.Addr, handlers map[string]func(conn net.Conn, params string)) {
	fmt.Println(addr.String(), "comes")
	t1 := time.Now()
	for {
		length_prefix := make([]byte, 4)
		n, _ := conn.Read(length_prefix)
		if n == 0 {
			elapsed := time.Since(t1)
			fmt.Println(addr.String(), "bye, use ", elapsed)
			conn.Close()
			return
		}
		length := binary.LittleEndian.Uint32(length_prefix)
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
		//fmt.Println(in_, params)
		handler := handlers[in_]
		handler(conn, params)
	}
}

func (s *Server) Loop() {
	for {
		conn, _ := s.Listener.Accept()
		go handler_conn(conn, conn.RemoteAddr(), s.Handlers) // 和v1的区别就是加了个go。。
	}
}

func ping(conn net.Conn, params string) {
	sendresult(conn, "pong", params)
}

func sendresult(conn net.Conn, out string, params string) {
	m := map[string]string{"out": out, "result": params}
	request, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("Fail to marshal, %s\n", err)
		return
	}
	length_prefix := make([]byte, 4)
	binary.LittleEndian.PutUint32(length_prefix, uint32(len(request)))
	conn.Write(length_prefix)
	conn.Write(request)
}
