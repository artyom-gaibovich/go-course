/*
Задание 1. Правильный тип ресивера.

Дана структура Counter с полем count int.
Реализуй два метода:
  - Inc()       — увеличивает счётчик на 1;
  - Value() int — возвращает текущее значение.

Подумай, какой ресивер (по значению или по указателю) нужен каждому методу.
После реализации программа должна напечатать 3.

Подсказка: метод, который МУТИРУЕТ получателя, обязан принимать его по указателю.
Не смешивай типы ресиверов между методами одного типа.
*/

package main

import "fmt"

type Counter struct {
	count int
}

// TODO: реализуй Inc()
func (c *Counter) Inc() {
	c.count++
}

// TODO: реализуй Value() int
func (c *Counter) Value() int {
	return c.count
}
func main() {
	c := &Counter{}
	c.Inc()
	c.Inc()
	c.Inc()
	fmt.Println(c.Value())
}
