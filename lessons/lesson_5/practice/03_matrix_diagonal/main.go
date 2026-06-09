/*
Задание 3: Главная диагональ матрицы

Реализуйте функцию MainDiagonal(matrix [][]int) []int, которая возвращает
главную диагональ квадратной матрицы (элементы matrix[i][i]).

Требования:
  - Матрица представлена как [][]int (слайс слайсов)
  - Если matrix == nil или len(matrix) == 0 — верните nil
  - Если матрица не квадратная (любая строка имеет длину ≠ len(matrix)) — верните nil
  - Используйте предвыделение capacity для результирующего слайса
  - Не создавайте вспомогательных аллокаций, кроме результирующего слайса

Примеры:
  MainDiagonal([][]int{{1,2},{3,4}})         → [1, 4]
  MainDiagonal([][]int{{1,2,3},{4,5,6},{7,8,9}}) → [1, 5, 9]
  MainDiagonal([][]int{{1,2},{3,4,5}})       → nil (не квадратная)
  MainDiagonal(nil)                          → nil

Дополнительно (по желанию):
  Реализуйте функцию AntiDiagonal(matrix [][]int) []int для побочной диагонали
  (элементы matrix[i][n-1-i]).
*/

package main

import "fmt"

func MainDiagonal(matrix [][]int) []int {
	// TODO: реализуйте функцию
	return nil
}

func AntiDiagonal(matrix [][]int) []int {
	// TODO: реализуйте функцию (дополнительно)
	return nil
}

func main() {
	cases := []struct {
		name   string
		matrix [][]int
		want   []int
	}{
		{
			"2x2",
			[][]int{{1, 2}, {3, 4}},
			[]int{1, 4},
		},
		{
			"3x3",
			[][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			[]int{1, 5, 9},
		},
		{
			"1x1",
			[][]int{{42}},
			[]int{42},
		},
		{
			"not square",
			[][]int{{1, 2}, {3, 4, 5}},
			nil,
		},
		{
			"nil",
			nil,
			nil,
		},
		{
			"empty",
			[][]int{},
			nil,
		},
	}

	fmt.Println("=== MainDiagonal ===")
	for _, tc := range cases {
		got := MainDiagonal(tc.matrix)
		ok := equalSlices(got, tc.want)
		fmt.Printf("%-15s → got=%v, want=%v, ok=%v\n", tc.name, got, tc.want, ok)
	}

	fmt.Println()
	fmt.Println("=== AntiDiagonal (дополнительно) ===")
	m := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	fmt.Println("AntiDiagonal 3x3:", AntiDiagonal(m))
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
