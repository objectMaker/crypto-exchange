package main

import (
	"fmt"
	"time"
)

type Order struct {
	Size      float64
	Limit     *Limit
	Bid       bool
	Timestamp int64
}

func (o *Order) String() string {
	return fmt.Sprintf("[size: %.2f]", o.Size)
}

type Limit struct {
	Price       float64
	Orders      []*Order
	TotalVolume float64
}

type Orderbook struct {
	Asks []*Limit
	Bids []*Limit
}

func NewLimit(price float64) *Limit {
	return &Limit{
		Price:  price,
		Orders: []*Order{},
	}
}

func NewOrder(bid bool, size float64) *Order {
	return &Order{
		Size:      size,
		Bid:       bid,
		Limit:     &Limit{},
		Timestamp: time.Now().UnixNano(),
	}
}

func (l *Limit) AddOrder(o *Order) {
	o.Limit = l
	l.Orders = append(l.Orders, o)
	l.TotalVolume += o.Size
}

func (l *Limit) DeleteOrder(o *Order) {
	orderLength := len(l.Orders)
	for i := 0; i < len(l.Orders); i++ {
		if l.Orders[i] == o {
			lastIndex := orderLength - 1
			// to remove current order
			l.Orders[i] = l.Orders[lastIndex]
			l.Orders = l.Orders[:lastIndex]
			break
		}
	}
	//to ensure garbage collector to work properly
	o.Limit = nil
	l.TotalVolume -= o.Size

}
