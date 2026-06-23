package main

import "fmt"

type Point struct {
	X int
	Y int
}

func main() {

	positional := Point{10, 20}

	named := Point{
		X: 10,
		Y: 20,
	}

	pointer := &Point{X: 10, Y: 20}

	manual := new(Point)
	manual.X = 10
	manual.Y = 20

	fmt.Println(positional, named, *pointer, *manual)
}
