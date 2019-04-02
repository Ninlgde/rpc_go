package main

import (
	"github.com/Ninlgde/rpc_go/v3.0"
	"github.com/gin-gonic/gin"
)

func main() {
	// init gin
	g := gin.Default()

	// init client
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

	g.Run("localhost:8888")
}
