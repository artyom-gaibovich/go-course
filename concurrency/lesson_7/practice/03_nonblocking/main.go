/*
Задание 3: Неблокирующие trySend / tryRecv.

Реализуйте две функции, которые НЕ блокируются:
  - trySend(ch chan<- int, v int) bool — пытается записать v, возвращает false,
    если записать прямо сейчас невозможно;
  - tryRecv(ch <-chan int) (int, bool) — пытается прочитать, возвращает ok=false,
    если читать прямо сейчас нечего.

В main проверьте поведение на буферизированном канале ёмкости 1: заполните его
и убедитесь, что второй trySend вернёт false, а tryRecv после опустошения — false.


*/

package main

import "fmt"

func trySend(ch chan<- int, v int) bool {
	select {
	case ch <- v:
		{
			return true
		}
	default:
		{
			return false
		}
	}
}

func tryRecv(ch <-chan int) (int, bool) {
	select {
	case v, ok := <-ch:
		if !ok {
			return 0, false
		}
		return v, true
	default:
		return 0, false
	}
}

func main() {
	ch := make(chan int, 1)
	r1 := trySend(ch, 1)
	r2 := trySend(ch, 1)
	fmt.Println(r1, r2)

	ch2 := make(chan int, 1)
	r3 := trySend(ch2, 1)
	v1, r1 := tryRecv(ch2)
	v2, r2 := tryRecv(ch2)

	fmt.Println(r1, r2, v1, v2, r3)
	// TODO: создать буферизированный канал и проверить trySend/tryRecv
}

//Подсказка: select { case ...: ...; default: ... }. Не используйте len/cap для
//проверки — это гонка.
