package main

import (
	"fmt"
	"testing"
)

func TestLimit(t *testing.T) {
	l := NewLimit(10_000)
	oa := NewOrder(true, 5)
	ob := NewOrder(true, 7)
	oc := NewOrder(true, 9)
	l.AddOrder(oa)
	l.AddOrder(ob)
	l.AddOrder(oc)
	l.DeleteOrder(oa)
	fmt.Println("new limit: ", l)
}

func TestOrderbook(t *testing.T) {
	ob := NewOrderBook()
	buyOrder := NewOrder(true, 10)
	buyOrder2 := NewOrder(true, 10)
	sellOrder2 := NewOrder(false, 10)

	ob.PlaceOrder(10_000, buyOrder)
	ob.PlaceOrder(10_000, buyOrder2)
	ob.PlaceOrder(10_000, sellOrder2)
	fmt.Printf("your ask: %v\n", ob.AskLimits[10_000])
	fmt.Printf("your bid: %v", ob.BidLimits[10_000])
}
