package main

import (
	"fmt"
	"unsafe"
)

// go build -gcflags='-m' . | grep escape

func main() {
	// умножение на 2^m
	var arrayWithStack [10 << 20]int8
	_ = arrayWithStack

	var arrayWithHeap [10<<20 + 1]int8
	_ = arrayWithHeap

	fmt.Println(unsafe.Sizeof(arrayWithHeap))

}
