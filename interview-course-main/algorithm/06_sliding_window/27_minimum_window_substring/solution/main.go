package main

import "fmt"

// Задача: Минимальное окно с подстрокой
//
// Решение: Скользящее окно с двумя указателями и частотными таблицами
//
// Идея:
// 1. Создаем частотную таблицу для строки t
// 2. Используем два указателя: left и right (границы окна)
// 3. Расширяем окно вправо (right++), пока не найдем все символы из t
// 4. Когда окно валидно (содержит все символы):
//    - Сохраняем его, если оно меньше предыдущего минимума
//    - Сжимаем окно слева (left++), пытаясь уменьшить
// 5. Повторяем пока right не достигнет конца
//
// Используем:
// - tCount: частоты символов в t
// - windowCount: частоты символов в текущем окне
// - required: количество уникальных символов в t
// - formed: количество символов с правильной частотой в окне
//
// Пример: s = "ADOBECODEBANC", t = "ABC"
// tCount = {A:1, B:1, C:1}, required = 3
//
// right=0: A, windowCount={A:1}, formed=1
// right=1: D, windowCount={A:1,D:1}, formed=1
// right=2: O, windowCount={A:1,D:1,O:1}, formed=1
// right=3: B, windowCount={A:1,D:1,O:1,B:1}, formed=2
// right=4: E, windowCount={A:1,D:1,O:1,B:1,E:1}, formed=2
// right=5: C, windowCount={A:1,D:1,O:1,B:1,E:1,C:1}, formed=3 → валидно! "ADOBEC"
// Сжимаем: left=1, удаляем A, formed=2, не валидно
// ...
// right=12: C, окно "BANC", formed=3 → валидно! меньше предыдущего
//
// Сложность:
// - Время: O(|s| + |t|) - каждый символ обрабатывается максимум дважды
// - Память: O(|s| + |t|) - для хранения частотных таблиц

func minWindow(s string, t string) string {
	if len(s) == 0 || len(t) == 0 || len(s) < len(t) {
		return ""
	}

	// Частотная таблица для t
	tCount := make(map[byte]int)
	for i := 0; i < len(t); i++ {
		tCount[t[i]]++
	}

	required := len(tCount) // количество уникальных символов в t
	formed := 0             // количество символов с правильной частотой

	// Частотная таблица для текущего окна
	windowCount := make(map[byte]int)

	// Результат: [длина окна, left, right]
	minLen := len(s) + 1
	minLeft, minRight := 0, 0

	left := 0

	// Расширяем окно вправо
	for right := 0; right < len(s); right++ {
		char := s[right]
		windowCount[char]++

		// Проверяем, достигли ли нужной частоты для этого символа
		if count, exists := tCount[char]; exists && windowCount[char] == count {
			formed++
		}

		// Пытаемся сжать окно слева, пока оно валидно
		for left <= right && formed == required {
			// Обновляем результат, если нашли меньшее окно
			if right-left+1 < minLen {
				minLen = right - left + 1
				minLeft = left
				minRight = right
			}

			// Удаляем символ слева
			leftChar := s[left]
			windowCount[leftChar]--
			if count, exists := tCount[leftChar]; exists && windowCount[leftChar] < count {
				formed--
			}

			left++
		}
	}

	if minLen == len(s)+1 {
		return ""
	}
	return s[minLeft : minRight+1]
}

func main() {
	// Тест 1: s = "ADOBECODEBANC", t = "ABC" → "BANC"
	s1, t1 := "ADOBECODEBANC", "ABC"
	result1 := minWindow(s1, t1)
	fmt.Printf("Test 1: s='%s', t='%s'\n", s1, t1)
	fmt.Printf("Result: '%s' (expected: 'BANC')\n\n", result1)

	// Тест 2: s = "a", t = "a" → "a"
	s2, t2 := "a", "a"
	result2 := minWindow(s2, t2)
	fmt.Printf("Test 2: s='%s', t='%s'\n", s2, t2)
	fmt.Printf("Result: '%s' (expected: 'a')\n\n", result2)

	// Тест 3: s = "a", t = "aa" → ""
	s3, t3 := "a", "aa"
	result3 := minWindow(s3, t3)
	fmt.Printf("Test 3: s='%s', t='%s'\n", s3, t3)
	fmt.Printf("Result: '%s' (expected: '')\n\n", result3)

	// Тест 4: s = "ab", t = "b" → "b"
	s4, t4 := "ab", "b"
	result4 := minWindow(s4, t4)
	fmt.Printf("Test 4: s='%s', t='%s'\n", s4, t4)
	fmt.Printf("Result: '%s' (expected: 'b')\n", result4)
}
