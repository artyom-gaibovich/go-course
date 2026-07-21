package main

import (
	"fmt"
	"math"
)

// Задача: Подмассив с заданной суммой
//
// Решение: Динамическое скользящее окно переменного размера
//
// Идея:
// 1. Используем два указателя: left и right (границы окна)
// 2. Расширяем окно вправо (right++), добавляя элементы
// 3. Когда сумма >= k:
//    - Сохраняем длину окна, если она минимальна
//    - Сжимаем окно слева (left++), удаляя элементы
//    - Повторяем, пока сумма >= k
// 4. Продолжаем до конца массива
//
// Пример: nums = [2,3,1,2,4,3], k = 7
// right=0: sum=2, sum<7
// right=1: sum=5, sum<7
// right=2: sum=6, sum<7
// right=3: sum=8, sum>=7! len=4 [2,3,1,2]
//   Сжимаем: left=1, sum=6, sum<7
// right=4: sum=10, sum>=7! len=3 [3,1,2,4]
//   Сжимаем: left=2, sum=7, sum>=7! len=3 [1,2,4]
//   Сжимаем: left=3, sum=6, sum<7
// right=5: sum=9, sum>=7! len=3 [2,4,3]
//   Сжимаем: left=4, sum=7, sum>=7! len=2 [4,3] (минимум!)
//   Сжимаем: left=5, sum=3, sum<7
// Минимальная длина = 2
//
// Важно:
// - Окно имеет переменный размер
// - Расширяем, когда сумма мала
// - Сжимаем, когда сумма достаточна
// - Работает только для положительных чисел
//
// Сложность:
// - Время: O(n) - каждый элемент добавляется и удаляется максимум один раз
// - Память: O(1) - используем только переменные

func minSubArrayLen(k int, nums []int) int {
	minLen := math.MaxInt32
	left := 0
	sum := 0

	// Расширяем окно вправо
	for right := 0; right < len(nums); right++ {
		sum += nums[right]

		// Сжимаем окно слева, пока сумма >= k
		for sum >= k {
			// Обновляем минимальную длину
			currentLen := right - left + 1
			if currentLen < minLen {
				minLen = currentLen
			}

			// Удаляем левый элемент
			sum -= nums[left]
			left++
		}
	}

	// Если не нашли подмассив
	if minLen == math.MaxInt32 {
		return 0
	}
	return minLen
}

func main() {
	// Тест 1: [2,3,1,2,4,3], k=7 → 2
	nums1 := []int{2, 3, 1, 2, 4, 3}
	result1 := minSubArrayLen(7, nums1)
	fmt.Printf("Test 1: nums=%v, k=7\n", nums1)
	fmt.Printf("Result: %d (expected: 2, subarray: [4,3])\n\n", result1)

	// Тест 2: [1,4,4], k=4 → 1
	nums2 := []int{1, 4, 4}
	result2 := minSubArrayLen(4, nums2)
	fmt.Printf("Test 2: nums=%v, k=4\n", nums2)
	fmt.Printf("Result: %d (expected: 1, subarray: [4])\n\n", result2)

	// Тест 3: [1,1,1,1,1,1,1,1], k=11 → 0
	nums3 := []int{1, 1, 1, 1, 1, 1, 1, 1}
	result3 := minSubArrayLen(11, nums3)
	fmt.Printf("Test 3: nums=%v, k=11\n", nums3)
	fmt.Printf("Result: %d (expected: 0, нет подмассива)\n\n", result3)

	// Тест 4: [1,2,3,4,5], k=11 → 3
	nums4 := []int{1, 2, 3, 4, 5}
	result4 := minSubArrayLen(11, nums4)
	fmt.Printf("Test 4: nums=%v, k=11\n", nums4)
	fmt.Printf("Result: %d (expected: 3, subarray: [3,4,5])\n", result4)
}
