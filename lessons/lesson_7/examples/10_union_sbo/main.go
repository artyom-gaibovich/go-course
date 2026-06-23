package main

import (
	"fmt"
	"unsafe"
)

type Buffer struct {
	size  int64
	union [16]byte
}

func main() {
	var small Buffer
	small.size = 5
	copy(small.union[:], []byte("hello"))
	fmt.Printf("small inline: %q\n", small.union[:small.size])

	var big Buffer
	big.size = 4096
	pointer := (*uintptr)(unsafe.Pointer(&big.union[0]))
	capacity := (*int64)(unsafe.Pointer(&big.union[8]))
	*pointer = 0xDEADBEEF
	*capacity = 4096
	fmt.Printf("big heap ptr: %#x cap: %d\n", *pointer, *capacity)
}
