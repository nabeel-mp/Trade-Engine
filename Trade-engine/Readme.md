High Performance Trade Engine (Go)

- In-memory order book
- Kafka event streaming
- Redis caching
- OHLCV candles (1s, 1m)
- Handles 200k+ orders/sec

Run:
docker-compose up
go run cmd/engine/main.go
go run cmd/candle/main.go
go run cmd/api/main.go
