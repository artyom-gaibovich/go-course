package main

import "fmt"

type Vector struct {
	X int
	Y int
	Z int
}

func wholeCopy(src Vector) Vector {
	return src
}

func fieldCopy(src Vector) Vector {
	var dst Vector
	dst.X = src.X
	dst.Y = src.Y
	dst.Z = src.Z
	return dst
}

func main() {
	src := Vector{X: 1, Y: 2, Z: 3}

	a := wholeCopy(src)
	b := fieldCopy(src)

	fmt.Println(a == b)

	src.X = 100
	fmt.Println(a.X, b.X)
}
