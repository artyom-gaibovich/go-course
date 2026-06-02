/*
Задание 1: Проверка чётности через битовую операцию

Реализуй функцию IsEven(n int) bool, которая определяет, является ли число чётным.
Используй побитовую операцию AND вместо оператора %.

Объясни (в комментарии): почему `n & 1 == 0` эквивалентно `n % 2 == 0`?

Дополнительно: напиши бенчмарк, сравнивающий два варианта — через & и через %.
Запусти: go test -bench=. main_test.go
*/

package main

import "fmt"

// IsEven возвращает true, если n чётное.
func IsEven(n int) bool {
	// TODO: реализовать через & без использования %
	return false
}

// IsOdd возвращает true, если n нечётное.
func IsOdd(n int) bool {
	// TODO: реализовать через IsEven
	return false
}

func main() {
	numbers := []int{0, 1, 2, -3, -4, 100, 101}
	for _, n := range numbers {
		fmt.Printf("%4d: even=%v, odd=%v\n", n, IsEven(n), IsOdd(n))
	}
}
