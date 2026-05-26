package main

import "fmt"

func main() {
	// Классический C-стиль: init; condition; post
	for i := 0; i < 5; i++ {
		fmt.Println(i)
	}

	// Обратный отсчёт
	for i := 10; i > 0; i -= 3 {
		fmt.Print(i, " ")
	}
	fmt.Println()

	// Несколько переменных в init и post
	for i, j := 0, 10; i < j; i, j = i+1, j-1 {
		fmt.Printf("i=%d j=%d\n", i, j)
	}
}
