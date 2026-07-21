package main

// Задача: Определить, являются ли две строки анаграммами
//
// Строка A — анаграмма строки B, если можно переставить местами символы
// в строке A и получить строку B.
//
// Требования:
// - Регистр символов учитывается
// - Пробелы учитываются
// - Пустая строка — анаграмма пустой строки
//
// Примеры:
// IsAnagram("лапоть", "пальто") → true
// IsAnagram("listen", "silent") → true
// IsAnagram("hello", "world") → false
// IsAnagram("", "") → true
// IsAnagram("a", "b") → false

func IsAnagram(a, b string) bool {
	// TODO: реализуйте функцию
	aMap := make(map[byte]int)
	bMap := make(map[byte]int)

	isAnagram := true
	for i := 0; i < len(a); i++ {
		aMap[a[i]]++
	}
	for i := 0; i < len(b); i++ {
		bMap[b[i]]++
	}

	if len(a) != len(b) {
		isAnagram = false
	}
	for k, v := range aMap {
		elemB, ok := bMap[k]
		if !ok {
			isAnagram = false
		}
		if elemB != v {
			isAnagram = false
		}

	}

	return isAnagram
}

func main() {
	testCases := []struct {
		a, b     string
		expected bool
	}{
		{"лапоть", "пальто", true},
		{"listen", "silent", true},
		{"hello", "world", false},
		{"", "", true},
		{"a", "b", false},
		{"анаграмма", "амаргана", false},
	}

	for _, tc := range testCases {
		result := IsAnagram(tc.a, tc.b)
		status := "✅"
		if result != tc.expected {
			status = "❌"
		}
		println(status, "IsAnagram('"+tc.a+"', '"+tc.b+"') =", result)
	}
}
