package main

import (
	"errors"
	"fmt"
)

// Go позволяет вернуть сразу несколько значений
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("деление на ноль")
	}
	return a / b, nil
}

// Идиома: (результат, ошибка) — самый распространённый паттерн в Go
func parseInt(s string) (int, error) {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	if err != nil {
		return 0, fmt.Errorf("не удалось распарсить %q: %w", s, err)
	}
	return n, nil
}

func main() {
	result, err := divide(10, 3)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Printf("10 / 3 = %.4f\n", result)
	}

	_, err = divide(5, 0)
	if err != nil {
		fmt.Println("Ошибка:", err)
	}

	if n, err := parseInt("42"); err == nil {
		fmt.Println("Число:", n)
	}
}
