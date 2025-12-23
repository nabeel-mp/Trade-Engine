package matching

import (
	"container/list"
	"sync"
)

type OrderBook struct {
	Bids map[float64]*list.List
	Asks map[float64]*list.List
	mu   sync.Mutex
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		Bids: make(map[float64]*list.List),
		Asks: make(map[float64]*list.List),
	}
}

func (ob *OrderBook) Add(order *Order) {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	book := ob.Bids
	if order.Side == Sell {
		book = ob.Asks
	}

	if _, ok := book[order.Price]; !ok {
		book[order.Price] = list.New()
	}
	book[order.Price].PushBack(order)
}
