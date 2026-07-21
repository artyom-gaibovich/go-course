package main

// Задача: Минимальное окно с подстрокой
//
// Даны две строки s и t. Необходимо найти минимальную подстроку в s,
// которая содержит все символы из t (включая дубликаты).
//
// Требования:
// - Сложность по времени O(n + m)
// - Использовать технику скользящего окна
// - Если такой подстроки нет, вернуть пустую строку
//
// Примеры:
// Input: s = "ADOBECODEBANC", t = "ABC"
// Output: "BANC"
// Объяснение: минимальная подстрока, содержащая A, B, C - это "BANC"
//
// Input: s = "a", t = "a"
// Output: "a"
//
// Input: s = "a", t = "aa"
// Output: ""
// Объяснение: в s нет двух символов 'a'

func minWindow(s string, t string) string {
	// TODO: реализуйте функцию
	return ""
}

func main() {
	// Тест 1
	result1 := minWindow("ADOBECODEBANC", "ABC")
	println("Test 1:", result1) // Ожидается: "BANC"

	// Тест 2
	result2 := minWindow("a", "a")
	println("Test 2:", result2) // Ожидается: "a"

	// Тест 3
	result3 := minWindow("a", "aa")
	println("Test 3:", result3) // Ожидается: ""
}
