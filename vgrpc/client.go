package vgrpc

import (
	"google.golang.org/grpc"
	"log"

	"context"
	"fmt"
	pb "github.com/Ninlgde/rpc_go/vgrpc/vgrpc"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"google.golang.org/grpc/connectivity"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

type GrpcClient struct {
	c       pb.GrpcServiceClient
	conn    *grpc.ClientConn
	address string
}

type Discover struct {
	dir   string
	nodes map[string]*GrpcClient
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
			discover.NewGrpcClient(string(kv.Key), string(kv.Value))
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
				discover.NewGrpcClient(string(event.Kv.Key), string(event.Kv.Value))
			case mvccpb.DELETE:
				discover.RemoveGrpcClient(string(event.Kv.Key))
				fmt.Println("DELETE event")
			}
			discover.Unlock()
		}
	}
}

func (discover *Discover) NewGrpcClient(key string, address string) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := pb.NewGrpcServiceClient(conn)
	client := &GrpcClient{
		c,
		conn,
		address,
	}

	discover.nodes[key] = client
}

func (discover *Discover) RemoveGrpcClient(key string) {
	client := discover.nodes[key]
	client.conn.Close()
	delete(discover.nodes, key)
}

type retryError struct{ error }

func (discover *Discover) FindAvailableClient() (pb.GrpcServiceClient, error) {
	for retry := 5; retry > 0; retry-- {
		discover.RLock()
		names := discover.NodeNames()
		node := names[rand.Int()%len(names)]
		rpcclient := discover.nodes[node]
		discover.RUnlock()
		c := rpcclient.c
		if rpcclient.conn.GetState() == connectivity.Shutdown || rpcclient.conn.GetState() == connectivity.TransientFailure {
			continue
		}
		return c, nil
	}
	return nil, retryError{}
}

type Client struct {
	discover *Discover
}

func NewClient(watch string) (client *Client) {
	dis := &Discover{
		dir:   watch,
		nodes: make(map[string]*GrpcClient),
	}
	go dis.Watch()

	//for {
	//	if len(dis.nodes) != 0 {
	//		break
	//	}
	//	time.Sleep(500 * time.Millisecond) // sleep 0.5 s
	//}
	client = &Client{dis}
	return
}

func (c *Client) Rpc(in_ string, params string) (out string, result string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	switch in_ {
	case "ping":
		gcli, err := c.discover.FindAvailableClient()
		if err != nil {
			fmt.Printf("Fail to conn, %s\n", err)
			return
		}
		r, err := gcli.PingCalc(ctx, &pb.PingRequest{Params: params})
		return r.Out, r.Result
	case "pi":
		gcli, err := c.discover.FindAvailableClient()
		if err != nil {
			fmt.Printf("Fail to conn, %s\n", err)
			return
		}
		n, _ := strconv.Atoi(params)
		r, err := gcli.PiCalc(ctx, &pb.PiRequest{N: int32(n)})
		return r.Out, strconv.FormatFloat(r.Value, 'E', -1, 64)
	}
	return
}
