package main

import "fmt"

// Задача: Максимальная сумма подмассива размера k
//
// Решение: Скользящее окно фиксированного размера
//
// Идея:
// 1. Вычисляем сумму первого окна размера k
// 2. Скользим окно вправо:
//    - Добавляем новый элемент справа
//    - Удаляем старый элемент слева
//    - Обновляем максимальную сумму
// 3. Возвращаем максимум
//
// Пример: nums = [2,1,5,1,3,2], k = 3
// Первое окно [2,1,5]: sum = 8
// Окно [1,5,1]: sum = 8 - 2 + 1 = 7
// Окно [5,1,3]: sum = 7 - 1 + 3 = 9 (максимум!)
// Окно [1,3,2]: sum = 9 - 5 + 2 = 6
// Максимальная сумма = 9
//
// Важно:
// - Не пересчитываем сумму окна каждый раз
// - Обновляем инкрементально: убираем левый элемент, добавляем правый
// - Это классическая техника скользящего окна
//
// Сложность:
// - Время: O(n) - один проход по массиву
// - Память: O(1) - используем только переменные

func maxSumSubarray(nums []int, k int) int {
	if len(nums) < k {
		return 0
	}

	// Вычисляем сумму первого окна
	windowSum := 0
	for i := 0; i < k; i++ {
		windowSum += nums[i]
	}

	maxSum := windowSum

	// Скользим окно вправо
	for i := k; i < len(nums); i++ {
		// Добавляем новый элемент справа, удаляем старый слева
		windowSum = windowSum + nums[i] - nums[i-k]

		// Обновляем максимум
		if windowSum > maxSum {
			maxSum = windowSum
		}
	}

	return maxSum
}

func main() {
	// Тест 1: [2,1,5,1,3,2], k=3 → 9
	nums1 := []int{2, 1, 5, 1, 3, 2}
	result1 := maxSumSubarray(nums1, 3)
	fmt.Printf("Test 1: nums=%v, k=3\n", nums1)
	fmt.Printf("Result: %d (expected: 9, subarray: [5,1,3])\n\n", result1)

	// Тест 2: [2,3,4,1,5], k=2 → 7
	nums2 := []int{2, 3, 4, 1, 5}
	result2 := maxSumSubarray(nums2, 2)
	fmt.Printf("Test 2: nums=%v, k=2\n", nums2)
	fmt.Printf("Result: %d (expected: 7, subarray: [3,4])\n\n", result2)

	// Тест 3: [1,4,2,10,23,3,1,0,20], k=4 → 39
	nums3 := []int{1, 4, 2, 10, 23, 3, 1, 0, 20}
	result3 := maxSumSubarray(nums3, 4)
	fmt.Printf("Test 3: nums=%v, k=4\n", nums3)
	fmt.Printf("Result: %d (expected: 39, subarray: [4,2,10,23])\n\n", result3)

	// Тест 4: [-1,-2,-3,-4], k=2 → -3
	nums4 := []int{-1, -2, -3, -4}
	result4 := maxSumSubarray(nums4, 2)
	fmt.Printf("Test 4: nums=%v, k=2\n", nums4)
	fmt.Printf("Result: %d (expected: -3, subarray: [-1,-2])\n", result4)
}
