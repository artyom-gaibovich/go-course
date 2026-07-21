package main

// Задача: Реверс слов (кроме палиндромов)
//
// Дана строка, содержащая слова, разделённые одиночными пробелами.
// Нужно написать функцию, которая вернёт строку, где каждое слово заменено на слово,
// состоящее из тех же букв, но в обратном порядке.
// Если слово является палиндромом (читается одинаково в обе стороны),
// его нужно оставить без изменений.
//
// Примеры:
// reverseWords("Hello worlD ollo") → "olleH Dlrow ollo"
// reverseWords("привет мир ара") → "тевирп рим ара"

func reverseWords(s string) string {
	// TODO: реализуйте функцию
	return ""
}

func main() {
	test1 := reverseWords("Hello worlD ollo")
	println("Test 1:", test1) // "olleH Dlrow ollo"

	test2 := reverseWords("привет мир ара")
	println("Test 2:", test2) // "тевирп рим ара"
}
