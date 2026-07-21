// Задача: Что выведет программа и почему?
// Объясните магию с capacity при создании срезов.

package main

import "fmt"

func main() {
	// Создаем слайс
	original := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Printf("original: len=%d cap=%d %v\n", len(original), cap(original), original)

	// Создаем срез от слайса
	slice1 := original[2:5]
	fmt.Printf("slice1 [2:5]: len=%d cap=%d %v\n", len(slice1), cap(slice1), slice1)

	// Еще один срез
	slice2 := original[3:7]
	fmt.Printf("slice2 [3:7]: len=%d cap=%d %v\n", len(slice2), cap(slice2), slice2)

	// Модифицируем slice1
	slice1[0] = 999
	fmt.Printf("После slice1[0]=999:\n")
	fmt.Printf("original: %v\n", original)
	fmt.Printf("slice1: %v\n", slice1)
	fmt.Printf("slice2: %v\n", slice2)

	// Append к slice1
	slice1 = append(slice1, 100, 200)
	fmt.Printf("\nПосле append(slice1, 100, 200):\n")
	fmt.Printf("slice1: len=%d cap=%d %v\n", len(slice1), cap(slice1), slice1)
	fmt.Printf("original: %v\n", original)
	fmt.Printf("slice2: %v\n", slice2)
}
