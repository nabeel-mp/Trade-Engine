package main

import (
	"Trade-engine/Internal/candle"
	"Trade-engine/Internal/kafka"
	"encoding/json"
	"time"
)

func main() {
	c1s := candle.Candle{}
	reader := kafka.NewConsumer("localhost:9092", "trades", "candle")

	kafka.Consume(reader, func(msg []byte) {
		var trade struct {
			Price    float64
			Quantity float64
		}
		_ = json.Unmarshal(msg, &trade)

		candle.Update(&c1s, trade.Price, trade.Quantity)

		if time.Since(c1s.Start) >= time.Second {
			c1s = candle.Candle{Start: time.Now()}
		}
	})
}
