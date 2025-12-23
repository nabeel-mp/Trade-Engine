package api

import (
	"Trade-engine/Internal/kafka"
	"Trade-engine/Internal/matching"
	"Trade-engine/Internal/redis"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var producer = kafka.NewProducer("localhost:9092", "orders")
var rdb = redis.New()

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
func getOrderBook(c *gin.Context) {
	val, err := rdb.Get(redis.Ctx, "orderbook").Result()
	if err != nil || val == "" {
		c.JSON(http.StatusOK, gin.H{"bids": gin.H{}, "asks": gin.H{}})
		return
	}
	var snapshot map[string]interface{}
	if err := json.Unmarshal([]byte(val), &snapshot); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse data"})
		return
	}
	c.JSON(http.StatusOK, snapshot)
}

// getTrades retrieves the last 100 executed trades
func getTrades(c *gin.Context) {
	trades, err := rdb.LRange(redis.Ctx, "trades", 0, 99).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch trades"})
		return
	}

	var result []interface{}
	for _, t := range trades {
		var mapped interface{}
		json.Unmarshal([]byte(t), &mapped)
		result = append(result, mapped)
	}
	c.JSON(http.StatusOK, result)
}

// getCandles retrieves the stored 1s OHLCV candles
func getCandles(c *gin.Context) {
	candles, err := rdb.LRange(redis.Ctx, "candles_1s", 0, -1).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch candles"})
		return
	}

	var result []interface{}
	for _, cand := range candles {
		var mapped interface{}
		json.Unmarshal([]byte(cand), &mapped)
		result = append(result, mapped)
	}
	c.JSON(http.StatusOK, result)
}
