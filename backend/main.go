package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/reserve/zoom/mtgs", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "All reserved Zoom MTG's!",
		})
	})
	r.POST("/reserve/zoom/mtg", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "reserved!",
		})
	})
	r.Run()
}
