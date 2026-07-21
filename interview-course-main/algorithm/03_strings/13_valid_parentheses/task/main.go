package main

// Задача: Валидация скобочной последовательности
//
// Дана строка, содержащая символы '(', ')', '{', '}', '[' и ']'.
// Необходимо определить, является ли данная строка валидной скобочной последовательностью.
//
// Скобочная последовательность является валидной при выполнении условий:
// 1. Открывающие скобки закрываются скобками того же типа
// 2. Открывающие скобки закрываются в правильном порядке
//
// Примеры:
// isValid("()") → true
// isValid("()[]{}") → true
// isValid("([])") → true
// isValid("([]}") → false (неправильный тип)
// isValid("({)}") → false (неправильный порядок)

func isValid(s string) bool {
	// TODO: реализуйте функцию
	return false
}

func main() {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"()", true},
		{"()[]{}", true},
		{"([])", true},
		{"([]}", false},
		{"({)}", false},
	}

	for _, tc := range testCases {
		result := isValid(tc.input)
		status := "✅"
		if result != tc.expected {
			status = "❌"
		}
		println(status, "isValid('"+tc.input+"') =", result)
	}
}
