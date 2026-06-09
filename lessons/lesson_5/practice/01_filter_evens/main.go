/*
Задание 1: Фильтрация чётных чисел

Реализуйте функцию FilterEvens(input []int) []int, которая возвращает новый
слайс, содержащий только чётные числа из входного слайса.

Требования:
  - Функция не должна модифицировать входной слайс
  - Предвыделите capacity для результата (сколько максимум элементов может быть?)
  - Обработайте nil-слайс на входе — верните nil
  - Обработайте пустой слайс — верните nil

Примеры:
  FilterEvens([]int{1, 2, 3, 4, 5, 6}) → [2, 4, 6]
  FilterEvens([]int{1, 3, 5})          → nil (нет чётных)
  FilterEvens(nil)                     → nil
  FilterEvens([]int{})                 → nil
*/

package main

import "fmt"

func FilterEvens(input []int) []int {
	// TODO: реализуйте функцию
	return nil
}

func main() {
	cases := []struct {
		name  string
		input []int
		want  []int
	}{
		{"mixed", []int{1, 2, 3, 4, 5, 6}, []int{2, 4, 6}},
		{"all odd", []int{1, 3, 5}, nil},
		{"all even", []int{2, 4, 6}, []int{2, 4, 6}},
		{"nil input", nil, nil},
		{"empty input", []int{}, nil},
		{"single even", []int{8}, []int{8}},
		{"single odd", []int{7}, nil},
	}

	for _, tc := range cases {
		got := FilterEvens(tc.input)
		ok := equalSlices(got, tc.want)
		fmt.Printf("%-15s → got=%v, want=%v, ok=%v\n", tc.name, got, tc.want, ok)
	}
}

func equalSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
