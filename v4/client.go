package v4_0

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"
)

type V4Client struct {
	conn    net.Conn
	address string
}

type Discover struct {
	dir   string
	nodes map[string]*V4Client
	sync.RWMutex
}

func (discover *Discover) NodeNames() (key []string) {
	discover.RLock()
	defer discover.RUnlock()
	key = make([]string, 0, len(discover.nodes))
	for k := range discover.nodes {
		key = append(key, k)
	}
	return
}

func (discover *Discover) Watch() {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		os.Exit(1)
	}

	var curRevision int64

	kv := clientv3.NewKV(client)
	for {
		rangeResp, err := kv.Get(context.TODO(), discover.dir, clientv3.WithPrefix())
		if err != nil {
			continue
		}

		discover.Lock()
		for _, kv := range rangeResp.Kvs {
			discover.NewV4Client(string(kv.Key), string(kv.Value))
		}
		discover.Unlock()

		// 从当前版本开始订阅
		curRevision = rangeResp.Header.Revision + 1
		break
	}

	watcher := clientv3.NewWatcher(client)
	watchChan := watcher.Watch(context.TODO(), discover.dir, clientv3.WithPrefix(), clientv3.WithRev(curRevision))
	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			discover.Lock()
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("PUT event")
				discover.NewV4Client(string(event.Kv.Key), string(event.Kv.Value))
			case mvccpb.DELETE:
				discover.RemoveV4Client(string(event.Kv.Key))
				fmt.Println("DELETE event")
			}
			discover.Unlock()
		}
	}
}

func (discover *Discover) NewV4Client(key string, address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := &V4Client{
		conn,
		address,
	}

	discover.nodes[key] = client
}

func (discover *Discover) RemoveV4Client(key string) {
	client := discover.nodes[key]
	client.conn.Close()
	delete(discover.nodes, key)
}

func (discover *Discover) FindAvailableConn() (net.Conn, error) {
	for retry := 5; retry > 0; retry-- {
		discover.RLock()
		names := discover.NodeNames()
		node := names[rand.Int()%len(names)]
		v4client := discover.nodes[node]
		discover.RUnlock()
		conn := v4client.conn
		// 假设都是活的
		return conn, nil
	}
	return nil, retryError{}
}

type Client struct {
	discover *Discover
}

type retryError struct{ error }

func (c *Client) Rpc(in_ string, params string) (out string, result string) {
	conn, err := c.discover.FindAvailableConn()
	if err != nil {
		fmt.Printf("Fail to conn, %s\n", err)
		return
	}
	m := map[string]string{"in": in_, "params": params}
	request, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("Fail to marshal, %s\n", err)
		return
	}
	length_prefix := make([]byte, 4)
	binary.LittleEndian.PutUint32(length_prefix, uint32(len(request)))
	conn.Write(length_prefix)
	conn.Write(request)
	conn.Read(length_prefix)
	length := binary.LittleEndian.Uint32(length_prefix)
	body := make([]byte, length)
	conn.Read(body)
	response := map[string]string{}
	err2 := json.Unmarshal(body, &response)
	if err2 != nil {
		fmt.Printf("Fail to unmarshal, %s\n", err2)
		return
	}
	return response["out"], response["result"]
}

func NewClient() (client *Client) {
	dis := &Discover{
		dir:   "ping",
		nodes: make(map[string]*V4Client),
	}
	go dis.Watch()

	for {
		if len(dis.nodes) != 0 {
			break
		}
		time.Sleep(500 * time.Millisecond) // sleep 0.5 s
	}
	client = &Client{dis}
	return
}
