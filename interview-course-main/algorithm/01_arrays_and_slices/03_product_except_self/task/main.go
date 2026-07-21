package main

// Задача: Произведение всех элементов кроме i-го
//
// Есть массив целых чисел. Нужно написать функцию, которая вернёт массив такого же размера,
// где на каждой i-ой позиции будет произведение всех элементов кроме i-го.
//
// Требование: Решить без использования деления
//
// Примеры:
// productExceptSelf([1, 2, 3]) → [2*3, 1*3, 1*2] → [6, 3, 2]
// productExceptSelf([1, 2, 3, 4]) → [24, 12, 8, 6]

func productExceptSelf(nums []int) []int {
	// TODO: реализуйте функцию
	return nil
}

func main() {
	result := productExceptSelf([]int{1, 2, 3})
	println("Test 1:", result) // [6, 3, 2]

	result = productExceptSelf([]int{1, 2, 3, 4})
	println("Test 2:", result) // [24, 12, 8, 6]
}
