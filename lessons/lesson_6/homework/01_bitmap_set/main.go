/*
Домашнее задание 1: Bitmap-множество

Реализуй тип BitSet для компактного хранения целых чисел из диапазона [0, 63].
Всё множество хранится в одном uint64 — каждый бит соответствует одному числу.

Требования:

  type BitSet uint64

  func (s *BitSet) Add(n int)          — добавить число n в множество
  func (s *BitSet) Remove(n int)       — удалить число n из множества
  func (s BitSet) Contains(n int) bool — вернуть true, если n есть в множестве
  func (s BitSet) Size() int           — количество элементов (используй PopCount, а не цикл по Contains)
  func (s BitSet) String() string      — вернуть строку вида "{1, 3, 7}" (числа отсортированы)

  Паника при n < 0 или n > 63.

Критерии:
  [ ] Add/Remove/Contains корректны на граничных значениях 0 и 63
  [ ] Size() реализован через подсчёт битов (алгоритм Кернигана или аналог)
  [ ] String() выводит числа в порядке возрастания через ", "
  [ ] Паника при выходе за границы диапазона

Пример работы:
  var s BitSet
  s.Add(1)
  s.Add(3)
  s.Add(63)
  s.Add(1)       // дубликат — не меняет состояние
  fmt.Println(s) // {1, 3, 63}
  fmt.Println(s.Size())       // 3
  fmt.Println(s.Contains(3))  // true
  s.Remove(3)
  fmt.Println(s.Contains(3))  // false
*/

package main

import "fmt"

type BitSet uint64

func (s *BitSet) Add(n int) {
	// TODO
}

func (s *BitSet) Remove(n int) {
	// TODO
}

func (s BitSet) Contains(n int) bool {
	// TODO
	return false
}

func (s BitSet) Size() int {
	// TODO: подсчёт через биты, без цикла по Contains
	return 0
}

func (s BitSet) String() string {
	// TODO: вернуть "{1, 3, 7}" или "{}" если пусто
	return "{}"
}

func main() {
	var s BitSet
	s.Add(1)
	s.Add(3)
	s.Add(63)
	s.Add(1)
	fmt.Println(s)
	fmt.Println("Size:", s.Size())
	fmt.Println("Contains(3):", s.Contains(3))
	s.Remove(3)
	fmt.Println("After Remove(3):", s)
	fmt.Println("Contains(3):", s.Contains(3))
}
