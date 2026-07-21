// Задача: Что выведет программа и почему?

package main

import "fmt"

func sliceCapacityDemo(data []int, start int) []int {
	subSlice := data[start:]
	return subSlice
}

func main() {
	data := make([]int, 5, 10)
	for i := range data {
		data[i] = i + 1
	}

	fmt.Printf("Изначальный слайс: %v, len %d, cap %d\n", data, len(data), cap(data))

	data = sliceCapacityDemo(data, 1)
	fmt.Printf("Подслайс (1): %v, len %d, cap %d\n", data, len(data), cap(data))

	dataNew := make([]int, 0, 3)
	dataNew = sliceCapacityDemo(data, 2)
	fmt.Printf("Подслайс (3): %v, len %d, cap %d\n", dataNew, len(dataNew), cap(dataNew))
}
