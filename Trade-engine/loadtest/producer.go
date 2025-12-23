package main

import (
	"Trade-engine/Internal/kafka"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

func main() {
	p := kafka.NewProducer("localhost:9092", "orders")
	jobs := make(chan int, 200000)
	var wg sync.WaitGroup

	// Start 100 workers
	for w := 1; w <= 100; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				order := map[string]interface{}{
					"id":        fmt.Sprintf("ord-%d", i),
					"price":     100.0 + (float64(i % 10)),
					"quantity":  1.0,
					"side":      "BUY",
					"type":      "LIMIT",
					"timestamp": time.Now(),
				}
				data, _ := json.Marshal(order)
				p.Publish(data)
			}
		}()
	}

	start := time.Now()
	for i := 0; i < 200000; i++ {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	fmt.Println("Sent 200k orders in:", time.Since(start))
}
