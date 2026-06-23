package main

import "fmt"

type Data struct {
	value int
}

func (d Data) Print() {
	fmt.Println("value:", d.value)
}

// d Data, аналог this - JavaScript, PHP, java, self = Python

func (Data) _() {}
func (Data) _() {}
func (Data) _() {}

func main() {
	d := Data{value: 42}
	d.Print()

	pointer := &Data{value: 7}
	pointer.Print()
}
