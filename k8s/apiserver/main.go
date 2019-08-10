package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	pb "github.com/Ninlgde/rpc_go/k8s/pb"
	"google.golang.org/grpc"
)

func main() {
	// Connect to GCD service
	conn, err := grpc.Dial("calc-service:3000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Dial failed: %v", err)
	}
	gcdClient := pb.NewGrpcServiceClient(conn)

	// Set up HTTP server
	r := gin.Default()
	r.GET("/gcd/:a/:b", func(c *gin.Context) {
		// Parse parameters
		a, err := strconv.ParseUint(c.Param("a"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameter A"})
			return
		}
		b, err := strconv.ParseUint(c.Param("b"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameter B"})
			return
		}
		// Call GCD service
		req := &pb.GCDRequest{A: a, B: b}
		if res, err := gcdClient.GCDCalc(c, req); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"result": fmt.Sprint(res.Result),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	})


	r.GET("/ping/:params", func(c *gin.Context) {
		p := c.Param("params")
		req := &pb.PingRequest{Params:p}
		if res, err := gcdClient.PingCalc(c, req); err == nil {

			c.JSON(http.StatusOK, gin.H{
				"out": res.Out,
				"result": res.Result,
			})
		}
	})
	r.GET("/pi/:n", func(c *gin.Context) {
		// Parse parameters
		n, err := strconv.Atoi(c.Param("n"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameter N"})
			return
		}
		req := &pb.PiRequest{N:int32(n)}
		if res, err := gcdClient.PiCalc(c, req); err == nil {

			c.JSON(http.StatusOK, gin.H{
				"out": res.Out,
				"result": res.Value,
			})
		}
	})

	// Run HTTP server
	if err := r.Run(":3000"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}