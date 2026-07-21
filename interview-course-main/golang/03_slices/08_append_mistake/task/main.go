// Задача: Что выведет программа и почему?
// Как правильно исправить appendLenWrong?

package main

import (
	"fmt"
)

func appendLenWrong(numbers []*int) {
	size := len(numbers)
	numbers = append(numbers, &size)
}

func main() {
	numbers := make([]*int, 0, 5)
	var number int
	for range 3 {
		number++
		numbers = append(numbers, &number)
	}

	appendLenWrong(numbers)

	for _, number := range numbers {
		fmt.Printf("%d ", *number)
	}
}
