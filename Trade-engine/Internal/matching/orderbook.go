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

// Add is the public method that handles locking
func (ob *OrderBook) Add(order *Order) {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	ob.addWithoutLock(order)
}

// addWithoutLock performs the actual insertion without touching the mutex
func (ob *OrderBook) addWithoutLock(order *Order) {
	book := ob.Bids
	if order.Side == Sell {
		book = ob.Asks
	}

	if _, ok := book[order.Price]; !ok {
		book[order.Price] = list.New()
	}
	book[order.Price].PushBack(order)
}
