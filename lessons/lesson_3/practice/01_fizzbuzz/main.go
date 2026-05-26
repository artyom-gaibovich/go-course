/*
Задание 1: FizzBuzz

Напечатай числа от 1 до 100, но:
  - вместо кратных 3 печатай "Fizz"
  - вместо кратных 5 печатай "Buzz"
  - вместо кратных и 3, и 5 — "FizzBuzz"

Пример вывода:
  1
  2
  Fizz
  4
  Buzz
  Fizz
  7
  ...
  FizzBuzz  (15, 30, 45, ...)

Подсказка: проверяй делимость на 15 первой.
*/

package main

import "fmt"

func main() {
	for i := 1; i <= 100; i++ {
		switch {
		case i%15 == 0:
			// TODO: напечатай "FizzBuzz"
		case i%3 == 0:
			// TODO: напечатай "Fizz"
		case i%5 == 0:
			// TODO: напечатай "Buzz"
		default:
			fmt.Println(i)
		}
	}
}
