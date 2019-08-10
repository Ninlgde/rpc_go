package main

import (
	"log"
	"net"

	pb "github.com/Ninlgde/rpc_go/k8s/pb"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"math"
)

type server struct{}

func main() {
	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGrpcServiceServer(s, &server{})
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *server) GCDCalc(ctx context.Context, in *pb.GCDRequest) (*pb.GCDResponse, error) {
	log.Printf("Received: %v : %v", in.A, in.B)
	a, b := in.A, in.B
	for b != 0 {
		a, b = b, a%b
	}
	return &pb.GCDResponse{Result: a}, nil
}

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