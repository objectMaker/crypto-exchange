package main

import (
	"fmt"
	"sort"
	"time"
)

type Match struct {
	Ask        *Order
	Bid        *Order
	SizeFilled float64
	Price      float64
}

type Order struct {
	Size      float64
	Limit     *Limit
	Bid       bool
	Timestamp int64
}

type Orders []*Order

func (o Orders) Len() int { return len(o) }

func (o Orders) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o Orders) Less(i, j int) bool { return o[i].Timestamp < o[j].Timestamp }

func (o *Order) String() string {
	return fmt.Sprintf("[size: %.2f]", o.Size)
}

type Limit struct {
	Price  float64
	Orders Orders
}

func (l *Limit) TotalVolume() float64 {
	var size float64
	for _, o := range l.Orders {
		size += o.Size
	}
	return size
}

func (l *Limit) Fill(o *Order) (matches []Match) {
	for _, order := range l.Orders {
		match := l.fillOrder(order, o)
		matches = append(matches, match)
		fmt.Println(order, "order++++")
		if order.Size == 0.0 {
			l.DeleteOrder(o)
		}
		if o.Size == 0.0 {
			break
		}
	}
	return
}

func (l *Limit) fillOrder(a, b *Order) Match {
	var (
		bid        *Order
		ask        *Order
		sizeFilled float64
	)

	if a.Bid {
		bid = a
		ask = b
	} else {
		bid = b
		ask = a
	}

	if a.Size >= b.Size {
		a.Size -= b.Size
		sizeFilled = b.Size
		b.Size = 0.0
	} else {
		b.Size -= a.Size
		sizeFilled = a.Size
		a.Size = 0.0
	}

	return Match{
		Bid:        bid,
		Ask:        ask,
		SizeFilled: sizeFilled,
		Price:      l.Price,
	}
}

type Limits []*Limit

type ByBestAsk struct{ Limits }

func (a ByBestAsk) Len() int {
	return len(a.Limits)
}
func (a ByBestAsk) Swap(i, j int) {
	a.Limits[i], a.Limits[j] = a.Limits[j], a.Limits[i]
}
func (a ByBestAsk) Less(i, j int) bool {
	return a.Limits[i].Price < a.Limits[j].Price
}

type ByBestBid struct{ Limits }

func (a ByBestBid) Len() int {
	return len(a.Limits)
}
func (a ByBestBid) Swap(i, j int) {
	a.Limits[i], a.Limits[j] = a.Limits[j], a.Limits[i]
}

// there have different symbol.
func (a ByBestBid) Less(i, j int) bool {
	return a.Limits[i].Price > a.Limits[j].Price
}

type Orderbook struct {
	asks      Limits
	bids      Limits
	AskLimits map[float64]*Limit
	BidLimits map[float64]*Limit
}

func (ob *Orderbook) TotalAskVolume() float64 {
	var TotalVolume float64
	for _, limit := range ob.asks {
		fmt.Println("LIMIT", limit)
		TotalVolume += limit.TotalVolume()
	}
	return TotalVolume
}

func (ob *Orderbook) TotalBidVolume() float64 {
	var TotalVolume float64
	for _, limit := range ob.bids {
		TotalVolume += limit.TotalVolume()
	}
	return TotalVolume
}

func NewOrderBook() *Orderbook {
	return &Orderbook{
		asks:      Limits{},
		bids:      Limits{},
		AskLimits: make(map[float64]*Limit),
		BidLimits: make(map[float64]*Limit),
	}
}

func NewLimit(price float64) *Limit {
	return &Limit{
		Price:  price,
		Orders: Orders{},
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
}

func (l *Limit) DeleteOrder(o *Order) {
	orderLength := len(l.Orders)
	lastIndex := orderLength - 1
	for i := 0; i < orderLength; i++ {
		if l.Orders[i] == o {
			// to remove current order
			l.Orders[i] = l.Orders[lastIndex]
			l.Orders = l.Orders[:lastIndex]
			break
		}
	}
	//to ensure garbage collector to work properly
	o.Limit = nil
}

func (ob *Orderbook) PlaceMarketOrder(o *Order) (matches []Match) {
	matches = []Match{}
	if o.Bid {
		fmt.Println("total ask", ob.TotalAskVolume())
		if o.Size > ob.TotalAskVolume() {
			panic(fmt.Sprintf("don't have enough asks,your bids is [%+v],total rest is [%+v]", o.Size, ob.TotalAskVolume()))
		}
		for _, limit := range ob.Asks() {
			limitMatches := limit.Fill(o)
			matches = append(matches, limitMatches...)

			if len(limit.Orders) == 0 {
				ob.clearLimit(true, limit)
			}
			if o.Size == 0.0 {
				break
			}
		}
	} else {
		if o.Size > ob.TotalBidVolume() {
			panic(fmt.Sprintf("don't have enough bids,your asks is [%+v],total rest is [%+v]", o.Size, ob.TotalBidVolume()))
		}
		for _, limit := range ob.Bids() {
			if limit.Price == o.Limit.Price {
				limitMatches := limit.Fill(o)
				matches = append(matches, limitMatches...)
			}
			if len(limit.Orders) == 0 {
				ob.clearLimit(true, limit)
			}
			if o.Size == 0.0 {
				break
			}
		}
	}
	return
}

func (ob *Orderbook) clearLimit(bid bool, limit *Limit) {
	if bid {

		for i := 0; i < len(ob.bids); i++ {
			if ob.bids[i] == limit {
				ob.bids[i] = ob.bids[len(ob.bids)-1]
				ob.bids = ob.bids[:len(ob.bids)-1]
			}
		}

	} else {
		for i := 0; i < len(ob.asks); i++ {
			if ob.asks[i] == limit {
				ob.asks[i] = ob.asks[len(ob.asks)-1]
				ob.asks = ob.asks[:len(ob.asks)-1]
			}
		}
	}
}

func (ob *Orderbook) PlaceLimitOrder(price float64, o *Order) {
	var limit *Limit
	if o.Bid {
		limit = ob.BidLimits[price]
	} else {
		limit = ob.AskLimits[price]
	}
	if limit == nil {
		limit = NewLimit(price)
		if o.Bid {
			ob.bids = append(ob.bids, limit)
			ob.BidLimits[price] = limit
		} else {
			ob.asks = append(ob.asks, limit)
			ob.AskLimits[price] = limit
		}
	}
	limit.AddOrder(o)
}

// sort the asks
func (ob *Orderbook) Asks() Limits {
	sort.Sort(ByBestAsk{ob.asks})
	return ob.asks
}

// sort the bids and return
func (ob *Orderbook) Bids() Limits {
	sort.Sort(ByBestBid{ob.bids})
	return ob.bids
}
