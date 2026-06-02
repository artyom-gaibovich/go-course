/*
Домашнее задание 2: Unique — дедупликация слайса

Реализуйте функцию:

    Unique[T comparable](s []T) []T

Которая возвращает слайс без дублирующихся элементов,
сохраняя порядок первого вхождения.

Требования:
  - Используйте map[T]struct{} для отслеживания виденных элементов.
  - Работайте in-place: result := s[:0], чтобы переиспользовать backing array.
  - В конце обнулите хвост: clear(s[len(result):]) — важно для GC.
  - Используйте дженерики (constraint: comparable).

Проверьте:
  Unique([]int{1, 2, 1, 3, 2, 4})           → [1 2 3 4]
  Unique([]string{"a", "b", "a", "c", "b"}) → [a b c]
  Unique([]int{})                            → []
  Unique([]int{5})                           → [5]

Bonus: объясните в комментарии — почему clear(s[len(result):]) важен
для слайсов, содержащих строки или указатели?
*/

package main

import "fmt"

// TODO: реализуйте Unique
func Unique[T comparable](s []T) []T {
	// TODO: создайте map seen
	// TODO: result := s[:0]
	// TODO: итерируйтесь по s, добавляйте в result если не видели
	// TODO: clear(s[len(result):]) для освобождения ссылок в хвосте
	// TODO: return result
	return nil
}

func main() {
	fmt.Println(Unique([]int{1, 2, 1, 3, 2, 4}))
	fmt.Println(Unique([]string{"a", "b", "a", "c", "b"}))
	fmt.Println(Unique([]int{}))
	fmt.Println(Unique([]int{5}))
}
