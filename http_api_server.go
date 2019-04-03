package main

import (
	"github.com/Ninlgde/rpc_go/v4"
	"github.com/Ninlgde/rpc_go/v5"
	"github.com/Ninlgde/rpc_go/vgrpc"
	"github.com/gin-gonic/gin"
)

func main() {
	// init gin
	g := gin.Default()

	//init v5 client
	client := v5_0.NewClient()

	g.GET("v5/ping/:params", func(c *gin.Context) {
		p := c.Param("params")
		out, params := client.Rpc("ping", p)
		c.JSON(200, gin.H{
			out: params,
		})
	})
	g.GET("v5/pi/:n", func(c *gin.Context) {
		n := c.Param("n")
		out, params := client.Rpc("pi", n)
		c.JSON(200, gin.H{
			out: params,
		})
	})

	//init v4 client
	v4client := v4_0.NewClient()

	g.GET("v4/ping/:params", func(c *gin.Context) {
		p := c.Param("params")
		out, params := v4client.Rpc("ping", p)
		c.JSON(200, gin.H{
			out: params,
		})
	})
	g.GET("v4/pi/:n", func(c *gin.Context) {
		n := c.Param("n")
		out, params := v4client.Rpc("pi", n)
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
