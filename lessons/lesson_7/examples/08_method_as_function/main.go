package main

import "fmt"

type Data struct {
	value int
}

func (d Data) Print() {
	fmt.Println("value:", d.value)
}

func main() {
	d := Data{value: 42}

	d.Print()

	Data.Print(d)

	method := Data.Print
	method(d)
}
