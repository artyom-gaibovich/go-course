/*
Задание 2: Дедупликация отсортированного слайса

Реализуйте функцию Dedup(s []int) []int, которая удаляет дублирующиеся
элементы из ОТСОРТИРОВАННОГО слайса, возвращая только уникальные значения.

Требования:
  - Решите задачу без дополнительной аллокации: работайте in-place
  - Используйте full slice expression s[:write:write] для возврата результата,
    чтобы caller не мог случайно получить доступ к «хвосту» через append
  - Обработайте nil и пустой слайс — верните nil
  - Порядок элементов должен сохраняться

Алгоритм:
  Держите «указатель записи» write (начиная с 1).
  Для каждого следующего элемента: если он отличается от s[write-1], запишите
  его в s[write] и увеличьте write.

Примеры:
  Dedup([]int{1, 1, 2, 3, 3, 3, 4}) → [1, 2, 3, 4]
  Dedup([]int{1, 2, 3})             → [1, 2, 3]
  Dedup([]int{5, 5, 5})             → [5]
  Dedup(nil)                        → nil
  Dedup([]int{})                    → nil
*/

package main

import "fmt"

func Dedup(s []int) []int {
	return nil
}

func main() {
	cases := []struct {
		name  string
		input []int
		want  []int
	}{
		{"with dupes", []int{1, 1, 2, 3, 3, 3, 4}, []int{1, 2, 3, 4}},
		{"no dupes", []int{1, 2, 3}, []int{1, 2, 3}},
		{"all same", []int{5, 5, 5}, []int{5}},
		{"two elements same", []int{3, 3}, []int{3}},
		{"nil", nil, nil},
		{"empty", []int{}, nil},
		{"single", []int{42}, []int{42}},
	}

	for _, tc := range cases {
		got := Dedup(tc.input)
		ok := equalSlices(got, tc.want)
		fmt.Printf("%-20s → got=%v, want=%v, ok=%v\n", tc.name, got, tc.want, ok)
	}

	fmt.Println()
	fmt.Println("Проверка: Dedup не должен давать caller доступ к 'хвосту' через append")
	original := []int{1, 1, 2, 2, 3}
	result := Dedup(original)
	if result != nil {
		fmt.Println("cap(result) == len(result)?", cap(result) == len(result))
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
