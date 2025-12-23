package api

import (
	"Trade-engine/Internal/kafka"
	"Trade-engine/Internal/matching"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var producer = kafka.NewProducer("localhost:9092", "orders")

func Register(r *gin.Engine) {
	r.POST("/order", placeOrder)
	r.GET("/orderbook", getOrderBook)
	r.GET("/trades", getTrades)
	r.GET("/candles", getCandles)
}

func placeOrder(c *gin.Context) {
	var order matching.Order

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if order.Timestamp.IsZero() {
		order.Timestamp = time.Now()
	}

	data, err := json.Marshal(order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "marshal failed"})
		return
	}

	if err := producer.Publish(data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "kafka publish failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "order published"})
}
