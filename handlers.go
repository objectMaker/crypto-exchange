package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PlaceOrderBody struct {
	UserID int64   `json:"userID"`
	Type   string  `json:"type"` //limit or market
	Bid    bool    `json:"bid"`
	Size   float64 `json:"size"`
	Price  float64 `json:"price"`
	Market string  `json:"market"`
}

func HandlePlaceOrder(ctx *gin.Context) {
	var placeOrderBody PlaceOrderBody
	ctx.Bind(&placeOrderBody)
	//business logic
	order := NewOrder(placeOrderBody.Bid, placeOrderBody.Size)
	if placeOrderBody.Type == "limit" {
		orderBook.PlaceLimitOrder(placeOrderBody.Price, order)
		ctx.JSON(http.StatusOK, gin.H{
			"message": map[string]any{
				"asks": orderBook.TotalAskVolume(),
				"bids": orderBook.TotalBidVolume(),
			},
		})
	} else {
		_, err := orderBook.PlaceMarketOrder(order)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"message": err.Error(),
			})
		}
	}

}
