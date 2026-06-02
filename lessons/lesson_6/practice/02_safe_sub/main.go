/*
Задание 2: Безопасное вычитание с защитой от переполнения

Реализуй функцию Sub(lhs, rhs int) (int, error), которая вычитает rhs из lhs
и возвращает ошибку ErrIntOverflow, если результат не вмещается в int.

Учти все комбинации знаков операндов. Для справки посмотри на реализацию Add
в lessons/lesson_6/examples/overflow_detection/main.go.

Подсказки:
  - lhs - rhs = lhs + (-rhs)
  - Будь осторожен: math.MinInt нельзя просто отрицать (это само по себе переполнение).
  - Когда rhs == math.MinInt, -rhs уже не помещается в int.

Проверь крайние случаи:
  Sub(math.MaxInt, -1)  → ошибка
  Sub(math.MinInt, 1)   → ошибка
  Sub(math.MinInt, math.MinInt) → 0, nil
  Sub(0, math.MinInt)   → ошибка
*/

package main

import (
	"errors"
	"fmt"
	"math"
)

var ErrIntOverflow = errors.New("integer overflow")

func Sub(lhs, rhs int) (int, error) {
	// TODO: реализовать
	return 0, nil
}

func main() {
	cases := [][2]int{
		{10, 3},
		{math.MaxInt, -1},
		{math.MinInt, 1},
		{math.MinInt, math.MinInt},
		{0, math.MinInt},
	}

	for _, c := range cases {
		result, err := Sub(c[0], c[1])
		if err != nil {
			fmt.Printf("Sub(%d, %d) = overflow\n", c[0], c[1])
		} else {
			fmt.Printf("Sub(%d, %d) = %d\n", c[0], c[1], result)
		}
	}
}
