package main

import "fmt"

type Data struct {
	value int
}

func (d Data) ByValue() {
	fmt.Println("by value:", d.value)
}

func (d *Data) ByPointer() {
	fmt.Println("by pointer:", d.value)
}

func main() {
	data := Data{value: 1}
	pointer := &data

	boundValue := pointer.ByValue
	boundPointer := pointer.ByPointer

	data.value = 100

	boundValue()
	boundPointer()
}
