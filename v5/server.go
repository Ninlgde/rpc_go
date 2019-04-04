package v5_0

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/Ninlgde/lrucache/go"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/etcdserver/api/v3rpc/rpctypes"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	Listener net.Listener
	Handlers map[string]func(conn net.Conn, params string)
}

var lruCache lru.LRUCache

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

func Ping(conn net.Conn, params string) {
	sendresult(conn, "pong", params)
}

func Pi(conn net.Conn, params string) {
	result := lruCache.Find(params)
	if result == nil {
		s := 0.0
		n, _ := strconv.Atoi(params)
		for i := 0; i <= n; i++ {
			s += 1.0 / (2*float64(i) + 1) / (2*float64(i) + 1)
		}
		r := math.Sqrt(8 * s)
		lruCache.Add(params, strconv.FormatFloat(r, 'E', -1, 64))
	}
	sendresult(conn, "pi_response", result.(string))
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

func Register(dir string, value string) {
	dir = strings.TrimRight(dir, "/") + "/"

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		os.Exit(1)
	}

	kv := clientv3.NewKV(client)
	lease := clientv3.NewLease(client)
	var curLeaseId clientv3.LeaseID = 0

	for {
		if curLeaseId == 0 {
			leaseResp, err := lease.Grant(context.TODO(), 10)
			if err != nil {
				goto SLEEP
			}

			key := dir + fmt.Sprintf("%d", leaseResp.ID)
			if _, err := kv.Put(context.TODO(), key, value, clientv3.WithLease(leaseResp.ID)); err != nil {
				goto SLEEP
			}
			curLeaseId = leaseResp.ID
		} else {
			//fmt.Printf("keepalive curLeaseId=%d\n", curLeaseId)
			if _, err := lease.KeepAliveOnce(context.TODO(), curLeaseId); err == rpctypes.ErrLeaseNotFound {
				curLeaseId = 0
				continue
			}
		}
	SLEEP:
		time.Sleep(time.Duration(1) * time.Second)
	}
}

func NewServer(address string) {
	// 全局lru
	lruCache = lru.NewLRUCache(10000)
	// 注册服务
	go Register("ping_lru", address)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("Fail to connect, %s\n", err)
	}
	defer listener.Close()
	handlers := make(map[string]func(conn net.Conn, params string))
	handlers["ping"] = Ping
	handlers["pi"] = Pi
	server := &Server{
		Listener: listener,
		Handlers: handlers,
	}
	server.Loop()
}
