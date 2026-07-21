// Задача: Что выведет программа и почему?

package main

import "fmt"

func main() {
	x := 10
	y := 20

	defer func(val int) {
		fmt.Println("x:", val)
	}(x)

	defer func() {
		fmt.Println("y:", y)
	}()

	x = 100
	y = 200

	fmt.Println("Конец main")
}
