package main

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
)

func main() {
	fmt.Println(uniqN(10))
}

// uniqN генерирует слайс из n уникальных случайных чисел.
func uniqN(n int) []int {
	// Используем map для отслеживания уже сгенерированных чисел.
	seen := make(map[int]bool)
	result := make([]int, 0, n)

	for len(result) < n {
		// Генерируем криптографически безопасное случайное число.
		num := cryptoRandInt(1000) // Диапазон можно настроить

		// Проверяем, не встречалось ли это число ранее.
		if !seen[num] {
			seen[num] = true
			result = append(result, num)
		}
	}

	return result
}

// cryptoRandInt генерирует криптографически безопасное случайное число от 0 до max-1.
func cryptoRandInt(max int) int {
	if max <= 0 {
		return 0
	}

	// Читаем случайные байты.
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	// Преобразуем байты в uint64.
	val := binary.BigEndian.Uint64(b)

	// Приводим к диапазону [0, max).
	return int(val % uint64(max))
}

// Альтернативное решение без map (если диапазон чисел ограничен):
// func uniqN(n int) []int {
//     max := n * 10 // Увеличиваем диапазон для уменьшения коллизий
//     used := make([]bool, max)
//     result := make([]int, 0, n)
//
//     for len(result) < n {
//         num := cryptoRandInt(max)
//         if !used[num] {
//             used[num] = true
//             result = append(result, num)
//         }
//     }
//     return result
// }

// Объяснение:
// 1. Используем crypto/rand вместо math/rand для криптографической безопасности.
// 2. Используем map для O(1) проверки уникальности.
// 3. Генерируем числа до тех пор, пока не получим n уникальных.
// 4. Предварительно выделяем память для result с capacity=n для оптимизации.
// 5. crypto/rand медленнее, но безопаснее для криптографических целей.
