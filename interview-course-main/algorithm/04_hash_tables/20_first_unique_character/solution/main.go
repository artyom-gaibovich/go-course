package main

import "fmt"

// Задача: Первый уникальный символ в строке
//
// Решение: Два прохода с хеш-таблицей
//
// Идея:
// 1. Первый проход: подсчитываем частоту каждого символа в map
// 2. Второй проход: ищем первый символ с частотой 1
//
// Пример: s = "loveleetcode"
// Первый проход - подсчет частот:
//   l: 2, o: 2, v: 1, e: 4, t: 1, c: 1, d: 1
//
// Второй проход - поиск первого уникального:
//   l (частота 2) - пропускаем
//   o (частота 2) - пропускаем
//   v (частота 1) - найден! индекс 2
//
// Альтернативный подход:
// Можно использовать массив размером 26 для английских букв,
// если нужна оптимизация по памяти для ограниченного алфавита.
//
// Сложность:
// - Время: O(n) - два прохода по строке
// - Память: O(1) или O(k) где k - размер алфавита (для английских букв k=26)

func firstUniqChar(s string) int {
	// Частотная таблица
	freq := make(map[rune]int)

	// Первый проход: подсчитываем частоты
	for _, char := range s {
		freq[char]++
	}

	// Второй проход: ищем первый символ с частотой 1
	for i, char := range s {
		if freq[char] == 1 {
			return i
		}
	}

	// Уникального символа не найдено
	return -1
}

// Альтернативное решение с массивом (для английских букв)
func firstUniqCharArray(s string) int {
	// Массив для подсчета частот (26 букв)
	freq := make([]int, 26)

	// Первый проход: подсчитываем частоты
	for _, char := range s {
		freq[char-'a']++
	}

	// Второй проход: ищем первый символ с частотой 1
	for i, char := range s {
		if freq[char-'a'] == 1 {
			return i
		}
	}

	return -1
}

func main() {
	// Тест 1: "leetcode" → 0 (символ 'l')
	s1 := "leetcode"
	fmt.Printf("Test 1: '%s' → %d (expected: 0, char: 'l')\n", s1, firstUniqChar(s1))

	// Тест 2: "loveleetcode" → 2 (символ 'v')
	s2 := "loveleetcode"
	fmt.Printf("Test 2: '%s' → %d (expected: 2, char: 'v')\n", s2, firstUniqChar(s2))

	// Тест 3: "aabb" → -1 (нет уникальных)
	s3 := "aabb"
	fmt.Printf("Test 3: '%s' → %d (expected: -1)\n", s3, firstUniqChar(s3))

	// Тест 4: "z" → 0 (единственный символ)
	s4 := "z"
	fmt.Printf("Test 4: '%s' → %d (expected: 0)\n", s4, firstUniqChar(s4))

	// Тест 5: альтернативное решение с массивом
	fmt.Println("\nАльтернативное решение с массивом:")
	fmt.Printf("Test 1 (array): '%s' → %d\n", s1, firstUniqCharArray(s1))
	fmt.Printf("Test 2 (array): '%s' → %d\n", s2, firstUniqCharArray(s2))
}
