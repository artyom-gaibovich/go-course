package main

import "fmt"

// Задача: Самая длинная подстрока без повторяющихся символов
//
// Решение: Скользящее окно (Sliding Window) с хеш-таблицей
//
// Идея:
// 1. Используем два указателя: left и right (границы окна)
// 2. Используем map для отслеживания последней позиции каждого символа
// 3. Расширяем окно вправо (right++)
// 4. Если встречаем повторяющийся символ:
//    - Сдвигаем левую границу за последнюю позицию этого символа
// 5. На каждом шаге обновляем максимальную длину
//
// Пример: s = "abcabcbb"
// Шаг 1: right=0, 'a', map={a:0}, окно="a", maxLen=1
// Шаг 2: right=1, 'b', map={a:0,b:1}, окно="ab", maxLen=2
// Шаг 3: right=2, 'c', map={a:0,b:1,c:2}, окно="abc", maxLen=3
// Шаг 4: right=3, 'a', 'a' уже есть на позиции 0, left=1, окно="bca", maxLen=3
// Шаг 5: right=4, 'b', 'b' уже есть на позиции 1, left=2, окно="cab", maxLen=3
// И так далее...
//
// Важно:
// - Храним индекс последнего вхождения символа + 1 для быстрого сдвига left
// - Левая граница может только увеличиваться (не уменьшается)
//
// Сложность:
// - Время: O(n) - каждый символ обрабатывается максимум дважды
// - Память: O(min(n, m)) где m - размер алфавита

func lengthOfLongestSubstring(s string) int {
	// Map для хранения последней позиции каждого символа
	charIndex := make(map[rune]int)
	maxLen := 0
	left := 0

	// Расширяем окно вправо
	for right, char := range s {
		// Если символ уже встречался и находится в текущем окне
		if lastIndex, exists := charIndex[char]; exists && lastIndex >= left {
			// Сдвигаем левую границу
			left = lastIndex + 1
		}

		// Обновляем позицию символа
		charIndex[char] = right

		// Вычисляем текущую длину окна
		currentLen := right - left + 1
		if currentLen > maxLen {
			maxLen = currentLen
		}
	}

	return maxLen
}

func main() {
	// Тест 1: "abcabcbb" → 3
	s1 := "abcabcbb"
	fmt.Printf("Test 1: '%s' → %d (expected: 3)\n", s1, lengthOfLongestSubstring(s1))

	// Тест 2: "bbbbb" → 1
	s2 := "bbbbb"
	fmt.Printf("Test 2: '%s' → %d (expected: 1)\n", s2, lengthOfLongestSubstring(s2))

	// Тест 3: "pwwkew" → 3
	s3 := "pwwkew"
	fmt.Printf("Test 3: '%s' → %d (expected: 3)\n", s3, lengthOfLongestSubstring(s3))

	// Тест 4: "" → 0
	s4 := ""
	fmt.Printf("Test 4: '%s' → %d (expected: 0)\n", s4, lengthOfLongestSubstring(s4))

	// Тест 5: "dvdf" → 3
	s5 := "dvdf"
	fmt.Printf("Test 5: '%s' → %d (expected: 3, substring: 'vdf')\n", s5, lengthOfLongestSubstring(s5))

	// Тест 6: "abba" → 2
	s6 := "abba"
	fmt.Printf("Test 6: '%s' → %d (expected: 2, substring: 'ab' or 'ba')\n", s6, lengthOfLongestSubstring(s6))
}
