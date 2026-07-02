package main

import (
	"fmt"
	"unsafe"
)

type Bad struct {
	a bool // 1
	// 3
	b int32 // 4

	c bool // 1
	// 3
}

type Good struct {
	b int32 // 4
	a bool  // 1
	c bool  // 1
	// 2
}

// struct{} - 0 байт

type FinalZero struct {
	x int64    // 8 байт
	z struct{} // 8 байт
}

type ZeroFirst struct {
	z struct{} // 0 байт
	x int64    // 8 байт
}

func main() {
	fmt.Println("Bad:", unsafe.Sizeof(Bad{}))
	fmt.Println("Good:", unsafe.Sizeof(Good{}))

	fmt.Println("FinalZero:", unsafe.Sizeof(FinalZero{}))
	fmt.Println("ZeroFirst:", unsafe.Sizeof(ZeroFirst{}))

	fmt.Println("Empty:", unsafe.Sizeof(struct{}{}))
}
