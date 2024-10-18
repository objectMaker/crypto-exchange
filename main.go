package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var orderBook *Orderbook

func main() {

	r := gin.Default()
	orderBook = NewOrderBook()
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "bilibili",
		})
	})
	r.POST("/order", HandlePlaceOrder)
	r.Run()
}
