package main

import "time"

type Order struct {
	Size      float64
	Limit     *Limit
	Bid       bool
	Timestamp int64
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
