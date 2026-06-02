/*
Задание 4: Стек на слайсе

Реализуйте generic-стек Stack[T any] на основе слайса:

    Push(v T)            — добавить элемент на вершину
    Pop() (T, bool)      — извлечь верхний элемент (false если стек пуст)
    Peek() (T, bool)     — посмотреть верхний элемент без удаления
    Len() int            — текущее количество элементов

Правила реализации:
  - Внутри — слайс []T.
  - Push → append в конец.
  - Pop → взять s[len-1], уменьшить len, обнулить освобождённый элемент.
  - При Pop/Peek пустого стека возвращайте zero value и false.

Проверьте:
  s.Push(1), s.Push(2), s.Push(3)
  s.Pop()  → 3, true
  s.Pop()  → 2, true
  s.Peek() → 1, true
  s.Pop()  → 1, true
  s.Pop()  → 0, false  (стек пуст)
  s.Len()  → 0
*/

package main

import "fmt"

// TODO: определите тип Stack[T any]

// TODO: Push
// TODO: Pop
// TODO: Peek
// TODO: Len

func main() {
	var s Stack[int]

	s.Push(1)
	s.Push(2)
	s.Push(3)

	fmt.Println(s.Pop())  // 3 true
	fmt.Println(s.Pop())  // 2 true
	fmt.Println(s.Peek()) // 1 true
	fmt.Println(s.Pop())  // 1 true
	fmt.Println(s.Pop())  // 0 false
	fmt.Println(s.Len())  // 0
}
