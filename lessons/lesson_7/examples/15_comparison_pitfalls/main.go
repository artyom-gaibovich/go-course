package main

import "fmt"

type Comparable struct {
	A int
	B string
}

func main() {
	x := Comparable{A: 1, B: "go"}
	y := Comparable{A: 1, B: "go"}
	fmt.Println("structs equal:", x == y)

	var i any = x
	var j any = y
	fmt.Println("via any equal:", i == j)

	withSlice := struct {
		data []int
	}{data: []int{1, 2}}

	var k any = withSlice
	fmt.Println("comparing slice-backed value panics:")
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered:", r)
		}
	}()
	fmt.Println(k == k)
}
