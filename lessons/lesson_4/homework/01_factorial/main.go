/*
Домашнее задание 1: Факториал

Реализуй две версии вычисления факториала n! = 1 * 2 * ... * n:

  factorialRecursive(n int) int  — рекурсивная версия
  factorialIterative(n int) int  — итеративная (через for)

Требования:
  - Для n = 0 → 1 (по определению)
  - Для n < 0 → паника с сообщением "факториал отрицательного числа не определён"
  - Обе версии должны возвращать одинаковые результаты

Проверь на:
  0! = 1
  1! = 1
  5! = 120
  10! = 3628800
  12! = 479001600  (максимум для int32, проверь тип!)

Подсказка по рекурсии: n! = n * (n-1)!
Усложнение: вычисли, при каком n результат переполняет int64.
*/

package main

import "fmt"

func factorialRecursive(n int) int {
	if n < 0 {
		panic("факториал отрицательного числа не определён")
	}
	// TODO
	return 0
}

func factorialIterative(n int) int {
	if n < 0 {
		panic("факториал отрицательного числа не определён")
	}
	// TODO
	return 0
}

func main() {
	tests := []int{0, 1, 5, 10, 12}
	for _, n := range tests {
		r := factorialRecursive(n)
		i := factorialIterative(n)
		match := "✓"
		if r != i {
			match = "✗ РАСХОЖДЕНИЕ"
		}
		fmt.Printf("%2d! = %12d  %s\n", n, r, match)
	}
}
