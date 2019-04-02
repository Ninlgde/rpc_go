package vgrpc

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"

	"fmt"
	pb "github.com/Ninlgde/rpc_go/vgrpc/vgrpc"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/etcdserver/api/v3rpc/rpctypes"
	"math"
	"os"
	"strings"
	"time"
)

type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) PingCalc(ctx context.Context, in *pb.PingRequest) (*pb.PingResponse, error) {
	log.Printf("Received: %v", in.Params)
	return &pb.PingResponse{Result: in.Params, Out: "pong"}, nil
}

func (s *server) PiCalc(ctx context.Context, in *pb.PiRequest) (*pb.PiResponse, error) {
	log.Printf("Received: %v", in.N)
	ss := 0.0
	n := int(in.N)
	for i := 0; i <= n; i++ {
		ss += 1.0 / (2*float64(i) + 1) / (2*float64(i) + 1)
	}
	result := math.Sqrt(8 * ss)
	return &pb.PiResponse{Value: result, Out: "pi_response"}, nil
}

func NewServer(address string) {
	go Register("grpc_ping", address)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGrpcServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
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
