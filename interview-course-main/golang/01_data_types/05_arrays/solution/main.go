package main

import "fmt"

// modifyArray принимает массив по значению (копия массива).
func modifyArray(arr [3]int) {
	// Изменение влияет только на локальную копию массива.
	arr[0] = 10
	fmt.Println("Inside modifyArray:", arr)
}

// modifySlice принимает слайс (который содержит указатель на массив).
func modifySlice(slice []int) {
	// Слайс содержит указатель на исходный массив, поэтому изменения видны снаружи.
	slice[0] = 5
	fmt.Println("Inside modifySlice:", slice)
}

func main() {
	array := [3]int{1, 2, 3}
	// Создаем слайс из массива. Слайс указывает на тот же массив.
	slice := array[:]

	fmt.Println("Before modifyArray:", array)
	modifyArray(array)                       // Передаем массив по значению (копия)
	fmt.Println("After modifyArray:", array) // Массив не изменился

	fmt.Println("Before modifySlice:", slice)
	modifySlice(slice)                       // Передаем слайс (содержит указатель на массив)
	fmt.Println("After modifySlice:", slice) // Слайс изменился
	fmt.Println("Final array:", array)       // Массив тоже изменился!
}

// Ответ:
// Before modifyArray: [1 2 3]
// Inside modifyArray: [10 2 3]
// After modifyArray: [1 2 3]
// Before modifySlice: [1 2 3]
// Inside modifySlice: [5 2 3]
// After modifySlice: [5 2 3]
// Final array: [5 2 3]
//
// Объяснение:
// 1. Массивы в Go передаются по значению (создается полная копия).
// 2. Слайсы передаются по значению, но слайс содержит указатель на массив.
//    Поэтому изменения элементов слайса видны в исходном массиве.
// 3. slice := array[:] создает слайс, который указывает на тот же массив.
// 4. Изменение через слайс изменяет исходный массив, так как они разделяют память.
