/*
Домашнее задание 1: Конкурентный стек.

Реализуйте потокобезопасный Stack[T any] на слайсе с КОРРЕКТНЫМ API:
  - Push(v T)
  - Pop() (v T, ok bool) — АТОМАРНО снимает и возвращает верхний элемент
  - Len() int

Требования:
  - Никаких раздельных Top()+Pop(): между ними может встроиться другая
    горутина и снять чужое значение (race на уровне API).
  - Указательный ресивер, мьютекс — приватное именованное поле.
  - Проверка на пустоту (ok == false при пустом стеке).
  - Всё под ОДНИМ мьютексом, Len тоже под блокировкой.

Проверка: запустите 100 consumer-горутин под `go run -race main.go`.
*/

package main

import (
	"fmt"
	"sync"
)

type Stack[T any] struct {
	mu   sync.Mutex
	data []T
}

func (s *Stack[T]) Push(v T) {
	// TODO
}

func (s *Stack[T]) Pop() (T, bool) {
	// TODO: под одной блокировкой проверить пустоту, снять и вернуть верхний
	var zero T
	return zero, false
}

func (s *Stack[T]) Len() int {
	// TODO
	return 0
}

func main() {
	var s Stack[int]
	for i := 0; i < 1000; i++ {
		s.Push(i)
	}

	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				s.Pop()
			}
		}()
	}
	wg.Wait()
	fmt.Println("len after 1000 pops:", s.Len()) // ожидается 0
}
