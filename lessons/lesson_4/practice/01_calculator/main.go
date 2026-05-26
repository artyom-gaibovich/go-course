/*
Задание 1: Безопасный калькулятор

Реализуй четыре функции:
  add(a, b float64) float64
  subtract(a, b float64) float64
  multiply(a, b float64) float64
  divide(a, b float64) (float64, error)  — ошибка при делении на 0

Затем напиши функцию calculate(a float64, op string, b float64) (float64, error),
которая принимает оператор в виде строки ("+", "-", "*", "/") и вызывает нужную функцию.
При неизвестном операторе возвращай ошибку.

Протестируй на наборе:
  10 + 3  → 13
  10 - 3  → 7
  10 * 3  → 30
  10 / 3  → 3.3333...
  10 / 0  → ошибка
  10 % 3  → ошибка (неизвестный оператор)
*/

package main

import (
	"fmt"
)

func add(a, b float64) float64 {
	return a + b
}

func subtract(a, b float64) float64 {
	// TODO
	return 0
}

func multiply(a, b float64) float64 {
	// TODO
	return 0
}

func divide(a, b float64) (float64, error) {
	// TODO
	return 0, nil
}

func calculate(a float64, op string, b float64) (float64, error) {
	// TODO: switch по op, вызывай нужную функцию
	return 0, fmt.Errorf("неизвестный оператор: %q", op)
}

func main() {
	tests := []struct {
		a, b float64
		op   string
	}{
		{10, 3, "+"},
		{10, 3, "-"},
		{10, 3, "*"},
		{10, 3, "/"},
		{10, 0, "/"},
		{10, 3, "%"},
	}

	for _, t := range tests {
		res, err := calculate(t.a, t.op, t.b)
		if err != nil {
			fmt.Printf("%.0f %s %.0f = ошибка: %v\n", t.a, t.op, t.b, err)
		} else {
			fmt.Printf("%.0f %s %.0f = %.4f\n", t.a, t.op, t.b, res)
		}
	}
}
