// Задача: Что выведет программа? Объясните почему.

package main

import "fmt"

type X struct {
	Val int
}

func main() {
	test()
}

func (x X) S() {
	fmt.Println(x.Val)
}

func test() {
	x := X{10}
	defer x.S()
	x.Val = 256
}
