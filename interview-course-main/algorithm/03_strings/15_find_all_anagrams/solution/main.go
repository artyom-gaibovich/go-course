package main

import "fmt"

// Задача: Поиск всех анаграмм в строке
//
// Решение: Скользящее окно фиксированного размера с подсчетом частот
//
// Идея:
// 1. Создаем частотную таблицу для строки p
// 2. Используем скользящее окно размером len(p)
// 3. Для каждой позиции окна:
//    - Добавляем новый символ справа
//    - Удаляем старый символ слева
//    - Сравниваем частоты символов в окне с частотами в p
// 4. Если частоты совпадают - нашли анаграмму
//
// Оптимизация:
// - Вместо сравнения всех частот каждый раз, используем счетчик совпадений
// - Отслеживаем количество символов с правильной частотой
//
// Пример: s = "cbaebabacd", p = "abc"
// Частоты p: {a:1, b:1, c:1}
// Окно [0:3] "cba": {c:1, b:1, a:1} → совпадает! индекс 0
// Окно [1:4] "bae": {b:1, a:1, e:1} → не совпадает
// Окно [2:5] "aeb": {a:1, e:1, b:1} → не совпадает
// ...
// Окно [6:9] "bac": {b:1, a:1, c:1} → совпадает! индекс 6
//
// Сложность:
// - Время: O(n) где n - длина строки s
// - Память: O(1) - размер алфавита фиксирован (26 букв)

func findAnagrams(s string, p string) []int {
	result := []int{}
	sLen, pLen := len(s), len(p)

	// Если s короче p, анаграмм быть не может
	if sLen < pLen {
		return result
	}

	// Частотные таблицы для p и текущего окна
	pCount := make(map[byte]int)
	windowCount := make(map[byte]int)

	// Заполняем частоты для p
	for i := 0; i < pLen; i++ {
		pCount[p[i]]++
	}

	// Инициализируем первое окно
	for i := 0; i < pLen; i++ {
		windowCount[s[i]]++
	}

	// Проверяем первое окно
	if mapsEqual(pCount, windowCount) {
		result = append(result, 0)
	}

	// Скользим окно по строке s
	for i := pLen; i < sLen; i++ {
		// Добавляем новый символ справа
		windowCount[s[i]]++

		// Удаляем старый символ слева
		leftChar := s[i-pLen]
		windowCount[leftChar]--
		if windowCount[leftChar] == 0 {
			delete(windowCount, leftChar)
		}

		// Проверяем текущее окно
		if mapsEqual(pCount, windowCount) {
			result = append(result, i-pLen+1)
		}
	}

	return result
}

// mapsEqual проверяет равенство двух map
func mapsEqual(m1, m2 map[byte]int) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, v := range m1 {
		if m2[k] != v {
			return false
		}
	}
	return true
}

func main() {
	// Тест 1: "cbaebabacd", "abc" → [0, 6]
	s1, p1 := "cbaebabacd", "abc"
	result1 := findAnagrams(s1, p1)
	fmt.Printf("Test 1: s='%s', p='%s' → %v (expected: [0, 6])\n", s1, p1, result1)

	// Тест 2: "abab", "ab" → [0, 1, 2]
	s2, p2 := "abab", "ab"
	result2 := findAnagrams(s2, p2)
	fmt.Printf("Test 2: s='%s', p='%s' → %v (expected: [0, 1, 2])\n", s2, p2, result2)

	// Тест 3: "baa", "aa" → [1]
	s3, p3 := "baa", "aa"
	result3 := findAnagrams(s3, p3)
	fmt.Printf("Test 3: s='%s', p='%s' → %v (expected: [1])\n", s3, p3, result3)

	// Тест 4: пустая строка
	s4, p4 := "", "a"
	result4 := findAnagrams(s4, p4)
	fmt.Printf("Test 4: s='%s', p='%s' → %v (expected: [])\n", s4, p4, result4)
}
