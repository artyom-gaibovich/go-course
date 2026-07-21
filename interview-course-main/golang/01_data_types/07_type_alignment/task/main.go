// Задача: Что выведет данный код? Объясните почему.

package main

import (
	"fmt"
	"unsafe"
)

type Foo struct {
	aaa bool
	bbb int32
	ccc bool
}

type Bar struct {
	aaa bool
	ccc bool
	bbb int32
}

func main() {
	ff := Foo{}
	bb := Bar{}
	fmt.Println(unsafe.Sizeof(ff))
	fmt.Println(unsafe.Sizeof(bb))
}
