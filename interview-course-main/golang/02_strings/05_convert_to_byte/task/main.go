// Задача: Как выполнить конвертацию из типа []byte в string без аллокаций?

package main

import (
	"fmt"
	"unsafe"
)

func main() {
	byteSlice := []byte("hello world")

	str1 := string(byteSlice)

	str2 := *(*string)(unsafe.Pointer(&byteSlice))

	fmt.Printf("%p\n", &byteSlice)
	fmt.Println(&str1)
	fmt.Println(&str2)

	fmt.Printf("%s\n", string(byteSlice))
	fmt.Println(str1)
	fmt.Println(str2)

	byteSlice[0] = 'R'

	fmt.Printf("%s\n", string(byteSlice))
	fmt.Println(str1)
	fmt.Println(str2)
}
