package main

import "fmt"

// Задача: Вращение массива
//
// Решение: Алгоритм с тройным разворотом
//
// Идея:
// 1. Нормализуем k (k = k % n), чтобы избежать лишних вращений
// 2. Развернуть весь массив
// 3. Развернуть первые k элементов
// 4. Развернуть оставшиеся n-k элементов
//
// Пример: nums = [1,2,3,4,5,6,7], k = 3
// Шаг 1: reverse весь массив -> [7,6,5,4,3,2,1]
// Шаг 2: reverse первые 3 элемента -> [5,6,7,4,3,2,1]
// Шаг 3: reverse оставшиеся 4 элемента -> [5,6,7,1,2,3,4]
//
// Сложность:
// - Время: O(n) - три прохода по массиву
// - Память: O(1) - используем только указатели

func rotate(nums []int, k int) {
	n := len(nums)
	if n == 0 {
		return
	}

	// Нормализуем k (если k больше длины массива)
	k = k % n
	if k == 0 {
		return
	}

	// Шаг 1: развернуть весь массив
	reverse(nums, 0, n-1)

	// Шаг 2: развернуть первые k элементов
	reverse(nums, 0, k-1)

	// Шаг 3: развернуть оставшиеся n-k элементов
	reverse(nums, k, n-1)
}

// reverse разворачивает элементы массива от start до end включительно
func reverse(nums []int, start, end int) {
	for start < end {
		nums[start], nums[end] = nums[end], nums[start]
		start++
		end--
	}
}

func main() {
	// Тест 1
	nums1 := []int{1, 2, 3, 4, 5, 6, 7}
	fmt.Println("Before:", nums1)
	rotate(nums1, 3)
	fmt.Println("After:", nums1) // Ожидается: [5,6,7,1,2,3,4]
	fmt.Println()

	// Тест 2
	nums2 := []int{-1, -100, 3, 99}
	fmt.Println("Before:", nums2)
	rotate(nums2, 2)
	fmt.Println("After:", nums2) // Ожидается: [3,99,-1,-100]
	fmt.Println()

	// Тест 3: k больше длины массива
	nums3 := []int{1, 2}
	fmt.Println("Before:", nums3)
	rotate(nums3, 3)             // 3 % 2 = 1, сдвиг на 1 позицию
	fmt.Println("After:", nums3) // Ожидается: [2,1]
}
