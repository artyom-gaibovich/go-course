/*
Задание 4: Счётчик через замыкание

Напиши функцию makeCounter() func() int, которая возвращает функцию-счётчик.
Каждый вызов счётчика увеличивает внутреннее состояние на 1 и возвращает новое значение.

Каждый вызов makeCounter() создаёт независимый счётчик.

Пример:
  counterA := makeCounter()
  counterB := makeCounter()

  fmt.Println(counterA()) // 1
  fmt.Println(counterA()) // 2
  fmt.Println(counterB()) // 1  ← независимый от A
  fmt.Println(counterA()) // 3

Усложнение: расширь до makeCounterWith(start, step int) func() int,
где start — начальное значение, step — шаг.
*/

package main

import "fmt"

func makeCounter() func() int {
	count := 0
	return func() int {
		// TODO: увеличь count и верни
		return count
	}
}

func makeCounterWith(start, step int) func() int {
	// TODO
	return func() int {
		return 0
	}
}

func main() {
	counterA := makeCounter()
	counterB := makeCounter()

	fmt.Println(counterA()) // 1
	fmt.Println(counterA()) // 2
	fmt.Println(counterB()) // 1
	fmt.Println(counterA()) // 3

	fmt.Println()
	byFive := makeCounterWith(0, 5)
	fmt.Println(byFive()) // 5
	fmt.Println(byFive()) // 10
	fmt.Println(byFive()) // 15
}
