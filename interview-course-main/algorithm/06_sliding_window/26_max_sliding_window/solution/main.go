package main

import "fmt"

// Задача: Максимум в скользящем окне
//
// Решение: Монотонная очередь (Deque - Double-Ended Queue)
//
// Идея:
// 1. Используем deque (очередь с двумя концами) для хранения индексов
// 2. Поддерживаем deque в убывающем порядке по значениям элементов
// 3. Первый элемент deque всегда содержит индекс максимума текущего окна
// 4. При добавлении нового элемента:
//    - Удаляем индексы за пределами окна (спереди)
//    - Удаляем индексы элементов меньше текущего (сзади)
//    - Добавляем текущий индекс в конец
//
// Пример: nums = [1,3,-1,-3,5,3,6,7], k = 3
// i=0: deque=[0], окно еще не готово
// i=1: nums[1]=3 > nums[0]=1, удаляем 0, deque=[1]
// i=2: nums[2]=-1 < nums[1]=3, добавляем, deque=[1,2], окно [1,3,-1], max=nums[1]=3
// i=3: удаляем индексы < 1, nums[3]=-3 < nums[2]=-1, deque=[1,2,3], окно [3,-1,-3], max=nums[1]=3
// i=4: удаляем индексы < 2, nums[4]=5 > все, deque=[4], окно [-1,-3,5], max=nums[4]=5
// i=5: nums[5]=3 < nums[4]=5, deque=[4,5], окно [-3,5,3], max=nums[4]=5
// i=6: удаляем индексы < 4, nums[6]=6 > nums[5]=3, deque=[6], окно [5,3,6], max=nums[6]=6
// i=7: nums[7]=7 > nums[6]=6, deque=[7], окно [3,6,7], max=nums[7]=7
//
// Важно:
// - Храним индексы, а не значения (для проверки границ окна)
// - Deque всегда в убывающем порядке
// - Операции добавления/удаления O(1) амортизированно
//
// Сложность:
// - Время: O(n) - каждый элемент добавляется и удаляется максимум один раз
// - Память: O(k) - размер deque не превышает k

func maxSlidingWindow(nums []int, k int) []int {
	if len(nums) == 0 || k == 0 {
		return []int{}
	}

	result := make([]int, 0, len(nums)-k+1)
	deque := make([]int, 0, k) // храним индексы

	for i := 0; i < len(nums); i++ {
		// Удаляем индексы за пределами окна (с начала)
		if len(deque) > 0 && deque[0] < i-k+1 {
			deque = deque[1:]
		}

		// Удаляем индексы элементов меньше текущего (с конца)
		// Поддерживаем убывающий порядок
		for len(deque) > 0 && nums[deque[len(deque)-1]] < nums[i] {
			deque = deque[:len(deque)-1]
		}

		// Добавляем текущий индекс
		deque = append(deque, i)

		// Когда окно заполнено, добавляем максимум в результат
		if i >= k-1 {
			result = append(result, nums[deque[0]])
		}
	}

	return result
}

func main() {
	// Тест 1: [1,3,-1,-3,5,3,6,7], k=3 → [3,3,5,5,6,7]
	nums1 := []int{1, 3, -1, -3, 5, 3, 6, 7}
	result1 := maxSlidingWindow(nums1, 3)
	fmt.Printf("Test 1: nums=%v, k=3\n", nums1)
	fmt.Printf("Result: %v (expected: [3,3,5,5,6,7])\n\n", result1)

	// Тест 2: [1], k=1 → [1]
	nums2 := []int{1}
	result2 := maxSlidingWindow(nums2, 1)
	fmt.Printf("Test 2: nums=%v, k=1\n", nums2)
	fmt.Printf("Result: %v (expected: [1])\n\n", result2)

	// Тест 3: [1,-1], k=1 → [1,-1]
	nums3 := []int{1, -1}
	result3 := maxSlidingWindow(nums3, 1)
	fmt.Printf("Test 3: nums=%v, k=1\n", nums3)
	fmt.Printf("Result: %v (expected: [1,-1])\n\n", result3)

	// Тест 4: [9,11], k=2 → [11]
	nums4 := []int{9, 11}
	result4 := maxSlidingWindow(nums4, 2)
	fmt.Printf("Test 4: nums=%v, k=2\n", nums4)
	fmt.Printf("Result: %v (expected: [11])\n", result4)
}
