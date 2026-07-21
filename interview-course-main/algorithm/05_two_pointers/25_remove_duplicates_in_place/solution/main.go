package main

import "fmt"

// Задача: Удаление дубликатов из отсортированного массива in-place
//
// Решение: Техника двух указателей (медленный и быстрый)
//
// Идея:
// 1. Используем два указателя:
//    - slow - указывает на позицию для записи уникального элемента
//    - fast - сканирует массив в поиске новых уникальных элементов
// 2. Начинаем с slow=0, fast=1 (первый элемент всегда уникален)
// 3. Если nums[fast] != nums[slow]:
//    - Нашли новый уникальный элемент
//    - Увеличиваем slow и копируем nums[fast] в nums[slow]
// 4. Всегда увеличиваем fast
// 5. Возвращаем slow+1 (количество уникальных элементов)
//
// Пример: nums = [0,0,1,1,1,2,2,3,3,4]
// slow=0, fast=1: nums[0]=0, nums[1]=0 (дубликат, пропускаем)
// slow=0, fast=2: nums[0]=0, nums[2]=1 (новый!) → slow=1, nums[1]=1
// slow=1, fast=3: nums[1]=1, nums[3]=1 (дубликат, пропускаем)
// slow=1, fast=4: nums[1]=1, nums[4]=1 (дубликат, пропускаем)
// slow=1, fast=5: nums[1]=1, nums[5]=2 (новый!) → slow=2, nums[2]=2
// slow=2, fast=6: nums[2]=2, nums[6]=2 (дубликат, пропускаем)
// slow=2, fast=7: nums[2]=2, nums[7]=3 (новый!) → slow=3, nums[3]=3
// slow=3, fast=8: nums[3]=3, nums[8]=3 (дубликат, пропускаем)
// slow=3, fast=9: nums[3]=3, nums[9]=4 (новый!) → slow=4, nums[4]=4
// Результат: [0,1,2,3,4], длина = 5
//
// Важно:
// - Массив отсортирован, поэтому дубликаты идут подряд
// - Не нужно проверять все элементы для поиска дубликатов
// - Изменяем массив на месте
//
// Сложность:
// - Время: O(n) - один проход по массиву
// - Память: O(1) - используем только два указателя

func removeDuplicates(nums []int) int {
	// Пустой массив или один элемент
	if len(nums) == 0 {
		return 0
	}

	// slow указывает на последний уникальный элемент
	slow := 0

	// fast сканирует массив
	for fast := 1; fast < len(nums); fast++ {
		// Если нашли новый уникальный элемент
		if nums[fast] != nums[slow] {
			slow++
			nums[slow] = nums[fast]
		}
	}

	// Возвращаем количество уникальных элементов
	return slow + 1
}

func main() {
	// Тест 1: [1,1,2] → 2, [1,2]
	nums1 := []int{1, 1, 2}
	fmt.Print("Before: ", nums1)
	k1 := removeDuplicates(nums1)
	fmt.Printf("\nAfter: k=%d, nums=%v (expected: k=2, nums=[1,2])\n\n", k1, nums1[:k1])

	// Тест 2: [0,0,1,1,1,2,2,3,3,4] → 5, [0,1,2,3,4]
	nums2 := []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}
	fmt.Print("Before: ", nums2)
	k2 := removeDuplicates(nums2)
	fmt.Printf("\nAfter: k=%d, nums=%v (expected: k=5, nums=[0,1,2,3,4])\n\n", k2, nums2[:k2])

	// Тест 3: [1,2,3] (нет дубликатов) → 3, [1,2,3]
	nums3 := []int{1, 2, 3}
	fmt.Print("Before: ", nums3)
	k3 := removeDuplicates(nums3)
	fmt.Printf("\nAfter: k=%d, nums=%v (expected: k=3, nums=[1,2,3])\n\n", k3, nums3[:k3])

	// Тест 4: [1] (один элемент) → 1, [1]
	nums4 := []int{1}
	fmt.Print("Before: ", nums4)
	k4 := removeDuplicates(nums4)
	fmt.Printf("\nAfter: k=%d, nums=%v (expected: k=1, nums=[1])\n", k4, nums4[:k4])
}
