package matching

import "time"

type Side string
type OrderType string

const (
	Buy  Side = "BUY"
	Sell Side = "SELL"

	Limit  OrderType = "LIMIT"
	Market OrderType = "MARKET"
)

type Order struct {
	ID        string    `json:"id"`
	Price     float64   `json:"price"`
	Quantity  float64   `json:"quantity"`
	Side      Side      `json:"side"`
	Type      OrderType `json:"type"`
	Timestamp time.Time `json:"timestamp"`
}
