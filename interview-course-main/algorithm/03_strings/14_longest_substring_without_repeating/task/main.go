package main

// Задача: Самая длинная подстрока без повторяющихся символов
//
// Дана строка s. Необходимо найти длину самой длинной подстроки без повторяющихся символов.
//
// Требования:
// - Сложность по времени O(n)
// - Использовать технику скользящего окна (sliding window)
//
// Примеры:
// Input: s = "abcabcbb"
// Output: 3
// Объяснение: подстрока "abc" имеет длину 3
//
// Input: s = "bbbbb"
// Output: 1
// Объяснение: подстрока "b" имеет длину 1
//
// Input: s = "pwwkew"
// Output: 3
// Объяснение: подстрока "wke" имеет длину 3

func lengthOfLongestSubstring(s string) int {
	// TODO: реализуйте функцию
	return 0
}

func main() {
	// Тест 1
	println("Test 1:", lengthOfLongestSubstring("abcabcbb")) // Ожидается: 3

	// Тест 2
	println("Test 2:", lengthOfLongestSubstring("bbbbb")) // Ожидается: 1

	// Тест 3
	println("Test 3:", lengthOfLongestSubstring("pwwkew")) // Ожидается: 3

	// Тест 4
	println("Test 4:", lengthOfLongestSubstring("")) // Ожидается: 0

	// Тест 5
	println("Test 5:", lengthOfLongestSubstring("dvdf")) // Ожидается: 3 (vdf)
}
