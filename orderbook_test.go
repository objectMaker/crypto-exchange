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

}
