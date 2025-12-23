package matching

import (
	"container/list"
	"sort"
	"sync"
)

type OrderBook struct {
	Bids      map[float64]*list.List
	Asks      map[float64]*list.List
	BidPrices []float64 // Keep sorted Descending
	AskPrices []float64 // Keep sorted Ascending
	mu        sync.Mutex
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		Bids:      make(map[float64]*list.List),
		Asks:      make(map[float64]*list.List),
		BidPrices: []float64{},
		AskPrices: []float64{},
	}
}

func (ob *OrderBook) addWithoutLock(order *Order) {
	book := ob.Bids
	prices := &ob.BidPrices
	isDescending := true

	if order.Side == Sell {
		book = ob.Asks
		prices = &ob.AskPrices
		isDescending = false
	}

	if _, ok := book[order.Price]; !ok {
		book[order.Price] = list.New()
		insertPrice(prices, order.Price, isDescending)
	}
	book[order.Price].PushBack(order)
}

func insertPrice(prices *[]float64, price float64, descending bool) {
	i := sort.Search(len(*prices), func(i int) bool {
		if descending {
			return (*prices)[i] <= price
		}
		return (*prices)[i] >= price
	})
	if i < len(*prices) && (*prices)[i] == price {
		return
	}
	*prices = append(*prices, 0)
	copy((*prices)[i+1:], (*prices)[i:])
	(*prices)[i] = price
}
