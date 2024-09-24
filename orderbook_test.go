package main

import "testing"

func TestLimit(t *testing.T) {
	l := NewLimit(10_000)
	l.AddOrder(NewOrder(true, 5))
	l.AddOrder(NewOrder(true, 7))
	l.AddOrder(NewOrder(true, 9))
	if l.TotalVolume != 21 {
		t.Error("not equal", l.TotalVolume)
	}
}

func TestOrderbook(t *testing.T) {

}
