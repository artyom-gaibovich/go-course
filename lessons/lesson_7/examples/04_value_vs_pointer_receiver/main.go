package main

import "fmt"

type Customer struct {
	balance int
}

func (c Customer) AddValue(v int) {
	c.balance += v
}

func (c *Customer) AddPointer(v int) {
	c.balance += v
}

func main() {
	c1 := Customer{}
	c1.AddValue(100)
	fmt.Println("value receiver:", c1.balance)

	c2 := Customer{}
	c2.AddPointer(100)
	fmt.Println("pointer receiver:", c2.balance)
}
