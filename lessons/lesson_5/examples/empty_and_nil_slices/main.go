package main

import (
	"fmt"
	"unsafe"
)

func main() {
	var data []string
	fmt.Println("var data []string:")
	fmt.Println(data)

	data = []string(nil)
	fmt.Println("data = []string(nil):")

	data = []string{}
	fmt.Println("data = []string{}:")

	data = make([]string, 10)
	data = append(data, "some str")
	fmt.Println("data = make([]string, 0):")

	empty := struct{}{}
	fmt.Println("empty struct address:", unsafe.Pointer(&empty))
}
