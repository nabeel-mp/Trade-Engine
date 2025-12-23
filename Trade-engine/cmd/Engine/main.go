package main

import (
	"encoding/json"
	"fmt"

	"Trade-engine/Internal/kafka"
	"Trade-engine/Internal/matching"
)

func main() {
	fmt.Println("🚀 Matching Engine started, waiting for orders...")

	ob := matching.NewOrderBook()

	reader := kafka.NewConsumer("localhost:9092", "orders", "engine")
	producer := kafka.NewProducer("localhost:9092", "trades")

	kafka.Consume(reader, func(msg []byte) {
		fmt.Println("📥 RAW MESSAGE:", string(msg))

		var order matching.Order
		if err := json.Unmarshal(msg, &order); err != nil {
			fmt.Println("❌ Unmarshal error:", err)
			return
		}

		fmt.Printf("🧾 Parsed Order: %+v\n", order)

		trades := ob.Match(&order)
		fmt.Println("🔍 Trades found:", len(trades))

		for _, t := range trades {
			data, _ := json.Marshal(t)
			_ = producer.Publish(data)
			fmt.Println("✅ Trade executed:", t)
		}
	})
}
