/*
Задание 2: Фильтрация слайса

Реализуйте функцию:

    Filter[T any](s []T, predicate func(T) bool) []T

Которая возвращает новый слайс только с теми элементами,
для которых predicate вернул true.

Требования:
  - Без лишних аллокаций: заранее выделите make([]T, 0, len(s)).
  - Не изменяйте исходный слайс.
  - Используйте дженерики (Go 1.21+).

Проверьте на []int{1, -2, 3, -4, 5} с предикатом > 0.
Ожидаемый результат: [1 3 5]

Дополнительно: проверьте на []string{"go", "rust", "c", "python"}
с предикатом len(s) > 2.
Ожидаемый результат: [rust python]
*/

package main

import "fmt"

// TODO: реализуйте функцию Filter
func Filter[T any](s []T, predicate func(T) bool) []T {
	// TODO
	return nil
}

func main() {
	nums := []int{1, -2, 3, -4, 5}
	positive := Filter(nums, func(n int) bool { return n > 0 })
	fmt.Println("positive:", positive)

	langs := []string{"go", "rust", "c", "python"}
	long := Filter(langs, func(s string) bool { return len(s) > 2 })
	fmt.Println("long names:", long)
}
