package main

import (
	"github.com/Ninlgde/rpc_go/v3.0"
	"github.com/Ninlgde/rpc_go/vgrpc"
	"github.com/gin-gonic/gin"
)

func main() {
	// init gin
	g := gin.Default()

	//init v3 client
	client := v3_0.NewClient()

	g.GET("v3/ping/:params", func(c *gin.Context) {
		p := c.Param("params")
		out, params := client.Rpc("ping", p)
		c.JSON(200, gin.H{
			out: params,
		})
	})
	g.GET("v3/pi/:n", func(c *gin.Context) {
		n := c.Param("n")
		out, params := client.Rpc("pi", n)
		c.JSON(200, gin.H{
			out: params,
		})
	})

	// init grpc client
	gclient := vgrpc.NewClient()

	g.GET("vgrpc/ping/:params", func(c *gin.Context) {
		p := c.Param("params")
		out, params := gclient.Rpc("ping", p)
		c.JSON(200, gin.H{
			out: params,
		})
	})
	g.GET("vgrpc/pi/:n", func(c *gin.Context) {
		n := c.Param("n")
		out, params := gclient.Rpc("pi", n)
		c.JSON(200, gin.H{
			out: params,
		})
	})

	g.Run("localhost:8888")
}
