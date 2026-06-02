package main

import "fmt"

func main() {
	var value int32 = 100
	var pointer *int32 = &value

	fmt.Println("address:", pointer)
	fmt.Println("value:", *pointer)
	fmt.Println("value without pointer:", value)

	*pointer = 500

	fmt.Println("address:", pointer)
	fmt.Println("value:", *pointer)
}
