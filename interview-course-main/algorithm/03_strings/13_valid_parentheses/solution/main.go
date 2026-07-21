package main

import "fmt"

// ============================================================================
// Решение: Стек
// ============================================================================

func isValid(s string) bool {
	var stack []rune
	pairs := map[rune]rune{
		')': '(',
		'}': '{',
		']': '[',
	}

	for _, char := range s {
		// Если закрывающая скобка
		if opening, isClosing := pairs[char]; isClosing {
			// Проверяем, что стек не пуст и вершина соответствует
			if len(stack) == 0 || stack[len(stack)-1] != opening {
				return false
			}
			// Удаляем из стека
			stack = stack[:len(stack)-1]
		} else {
			// Открывающая скобка — добавляем в стек
			stack = append(stack, char)
		}
	}

	// Валидная последовательность, если стек пуст
	return len(stack) == 0
}

func main() {
	fmt.Println("=== Решение: Стек ===")

	testCases := []struct {
		input    string
		expected bool
	}{
		{"()", true},
		{"()[]{}", true},
		{"([])", true},
		{"([}]", false},
		{"({)}", false},
		{"", true},
		{"(", false},
		{")", false},
	}

	for _, tc := range testCases {
		result := isValid(tc.input)
		status := "✅"
		if result != tc.expected {
			status = "❌"
		}
		fmt.Printf("%s isValid('%s') = %v (expected %v)\n",
			status, tc.input, result, tc.expected)
	}
}

// Временная сложность: O(n)
// - Один проход по строке, каждая операция со стеком O(1)
//
// Пространственная сложность: O(n)
// - В худшем случае все символы в стеке (например, "((((((")
