package main

import "fmt"

// Задача: Контейнер с максимальным количеством воды
//
// Решение: Техника двух указателей
//
// Идея:
// 1. Используем два указателя: left (начало) и right (конец)
// 2. Вычисляем площадь для текущей пары линий
// 3. Двигаем указатель с меньшей высотой внутрь
//    - Если height[left] < height[right], двигаем left++
//    - Иначе двигаем right--
// 4. Обновляем максимальную площадь
//
// Почему это работает:
// - Площадь = ширина * минимальная высота
// - Когда сдвигаем указатели, ширина уменьшается
// - Чтобы увеличить площадь, нужна большая высота
// - Двигаем указатель с меньшей высотой, т.к. он ограничивает площадь
// - Есть шанс найти более высокую линию
//
// Пример: height = [1,8,6,2,5,4,8,3,7]
// Шаг 1: left=0(1), right=8(7), area = 8 * min(1,7) = 8, двигаем left
// Шаг 2: left=1(8), right=8(7), area = 7 * min(8,7) = 49, двигаем right
// Шаг 3: left=1(8), right=7(3), area = 6 * min(8,3) = 18, двигаем right
// ...
// Максимальная площадь: 49
//
// Сложность:
// - Время: O(n) - один проход с двумя указателями
// - Память: O(1) - используем только переменные

func maxArea(height []int) int {
	maxArea := 0
	left := 0
	right := len(height) - 1

	for left < right {
		// Вычисляем текущую площадь
		width := right - left
		minHeight := min(height[left], height[right])
		currentArea := width * minHeight

		// Обновляем максимум
		if currentArea > maxArea {
			maxArea = currentArea
		}

		// Двигаем указатель с меньшей высотой
		if height[left] < height[right] {
			left++
		} else {
			right--
		}
	}

	return maxArea
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	// Тест 1: [1,8,6,2,5,4,8,3,7] → 49
	height1 := []int{1, 8, 6, 2, 5, 4, 8, 3, 7}
	fmt.Printf("Test 1: %v → %d (expected: 49)\n", height1, maxArea(height1))
	fmt.Println("Объяснение: между индексами 1(высота 8) и 8(высота 7): (8-1) * min(8,7) = 7 * 7 = 49")
	fmt.Println()

	// Тест 2: [1,1] → 1
	height2 := []int{1, 1}
	fmt.Printf("Test 2: %v → %d (expected: 1)\n", height2, maxArea(height2))
	fmt.Println()

	// Тест 3: [4,3,2,1,4] → 16
	height3 := []int{4, 3, 2, 1, 4}
	fmt.Printf("Test 3: %v → %d (expected: 16)\n", height3, maxArea(height3))
	fmt.Println("Объяснение: между индексами 0(высота 4) и 4(высота 4): (4-0) * min(4,4) = 4 * 4 = 16")
	fmt.Println()

	// Тест 4: [1,2,1] → 2
	height4 := []int{1, 2, 1}
	fmt.Printf("Test 4: %v → %d (expected: 2)\n", height4, maxArea(height4))
	fmt.Println("Объяснение: между индексами 0(высота 1) и 2(высота 1): (2-0) * min(1,1) = 2 * 1 = 2")
}
