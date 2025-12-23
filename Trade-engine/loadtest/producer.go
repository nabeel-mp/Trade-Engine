package main

import (
	"Trade-engine/Internal/kafka"
	"encoding/json"
	"fmt"
	"time"
)

func main() {
	p := kafka.NewProducer("localhost:9092", "orders")

	start := time.Now()
	for i := 0; i < 200000; i++ {
		order := map[string]interface{}{
			"id":        i,
			"price":     100.0,
			"quantity":  1.0,
			"side":      "BUY",
			"type":      "LIMIT",
			"timestamp": time.Now(),
		}
		data, _ := json.Marshal(order)
		go p.Publish(data)
	}
	fmt.Println("Sent in:", time.Since(start))
}
