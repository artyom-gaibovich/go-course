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

// Ответ:
// Изначальный слайс: [1 2 3 4 5], len 5, cap 10
// Подслайс (1): [2 3 4 5], len 4, cap 9
// Подслайс (3): [4 5], len 2, cap 7
//
// Объяснение:
// 1. data[start:] создает новый слайс, который начинается с индекса start.
// 2. Новый слайс указывает на тот же массив, что и исходный.
// 3. Длина нового слайса = len(data) - start.
// 4. Capacity нового слайса = cap(data) - start.
//    Это потому, что новый слайс "видит" только оставшуюся часть массива.
//
// Детали:
// - Изначально: data = [1 2 3 4 5], len=5, cap=10
// - После data[1:]: data = [2 3 4 5], len=4 (5-1), cap=9 (10-1)
// - После data[2:]: dataNew = [4 5], len=2 (4-2), cap=7 (9-2)
//
// Важно: dataNew имеет capacity=7, хотя мы создали его с make([]int, 0, 3).
// Это потому, что sliceCapacityDemo возвращает слайс, который указывает на
// тот же массив, что и data, и capacity определяется позицией в исходном массиве.
