package main

import (
	"Trade-engine/Internal/candle"
	"Trade-engine/Internal/kafka"
	"Trade-engine/Internal/redis" // Import internal redis package
	"encoding/json"
	"time"
)

func main() {
	rdb := redis.New() // Initialize redis client
	c1s := candle.Candle{Start: time.Now()}
	reader := kafka.NewConsumer("localhost:9092", "trades", "candle")

	kafka.Consume(reader, func(msg []byte) {
		var trade struct {
			Price    float64
			Quantity float64
		}
		_ = json.Unmarshal(msg, &trade)

		candle.Update(&c1s, trade.Price, trade.Quantity)

		if time.Since(c1s.Start) >= time.Second {
			data, _ := json.Marshal(c1s)
			rdb.LPush(redis.Ctx, "candles_1s", data)
			rdb.LTrim(redis.Ctx, "candles_1s", 0, 499)
			c1s = candle.Candle{Start: time.Now()}
		}
	})
}
