package main

import (
	"fmt"
	"reflect"
	"testing"
)

func assertEq(t *testing.T, a, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("%+v not equals %+v", a, b)
	}
}

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

func TestPlaceLimitOrder(t *testing.T) {
	ob := NewOrderBook()
	buyOrder := NewOrder(true, 10)
	buyOrder2 := NewOrder(true, 10)
	sellOrder2 := NewOrder(false, 10)

	ob.PlaceLimitOrder(10_000, buyOrder)
	ob.PlaceLimitOrder(10_001, buyOrder2)
	ob.PlaceLimitOrder(10_000, sellOrder2)
	fmt.Printf("your ask: %v\n", ob.AskLimits[10_000])
	fmt.Printf("your bid: %v", ob.BidLimits[10_000])
	//bid buy
	assertEq(t, len(ob.Bids()), 2)
	assertEq(t, len(ob.Asks()), 1)
}

func TestPlaceMarketOrder(t *testing.T) {
	ob := NewOrderBook()
	buyOrder1 := NewOrder(false, 10)
	buyOrder2 := NewOrder(false, 20)

	ob.PlaceLimitOrder(100, buyOrder1)
	ob.PlaceLimitOrder(100, buyOrder2)
	marketOrder := NewOrder(true, 15)
	matches := ob.PlaceMarketOrder(marketOrder)
	fmt.Println(matches, "matches")
}
