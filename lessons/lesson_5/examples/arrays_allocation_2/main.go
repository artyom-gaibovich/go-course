package main

import "fmt"

//go:noinline
func allocation() *[10]int {
	var data [10]int
	fmt.Println(&data)
	return &data
}

func main() {
	_ = allocation()
}
