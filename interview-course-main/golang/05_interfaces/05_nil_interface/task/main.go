// Задача: Что выведет программа и почему?

package main

import "fmt"

type SomeStruct struct {
	Value int
}

func CheckForNil(i interface{}) {
	if i == nil {
		fmt.Println("Это nil!")
		return
	}

	fmt.Println("Это не nil!")
}

func main() {
	var s *SomeStruct
	CheckForNil(s)
}
