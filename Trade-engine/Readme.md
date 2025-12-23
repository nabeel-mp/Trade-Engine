High Performance Trade Engine (Go)

- In-memory order book
- Kafka event streaming
- Redis caching
- OHLCV candles (1s, 1m)
- Handles 200k orders/sec

Run:
docker-compose up
go run cmd/engine/main.go
go run cmd/candle/main.go
go run cmd/api/main.go

## References
- Apache Kafka Documentation – Consumer groups, offsets, blocking consumption
- segmentio/kafka-go – High-performance Go Kafka client
- Redis Documentation – In-memory caching and fast read access
- Financial Market Concepts – Price-Time Priority matching
- Go Effective Go – Concurrency and synchronization
- Gin Web Framework – Typed JSON binding and REST APIs

