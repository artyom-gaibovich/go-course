/*
Задание 3: Трассировка функций через defer

Напиши вспомогательную функцию trace(name string) func(), которая:
  1. При вызове печатает: "→ вход: <name>"
  2. Возвращает функцию-очиститель, которая печатает: "← выход: <name>"

Использование:
  func foo() {
      defer trace("foo")()   // ← вызов trace сразу, defer откладывает возврат
      // ... тело функции
  }

Реализуй три взаимосвязанные функции: a() вызывает b(), b() вызывает c().
Каждая использует trace. Запусти и убедись в порядке вывода.

Ожидаемый вывод:
  → вход: a
  → вход: b
  → вход: c
  ← выход: c
  ← выход: b
  ← выход: a

Усложнение: добавь в trace вывод времени выполнения функции (time.Since).
*/

package main

import "fmt"

func trace(name string) func() {
	fmt.Println("→ вход:", name)
	return func() {
		// TODO: напечатай "← выход: <name>"
	}
}

func c() {
	defer trace("c")()
	fmt.Println("  тело c")
}

func b() {
	defer trace("b")()
	fmt.Println("  тело b")
	c()
}

func a() {
	defer trace("a")()
	fmt.Println("  тело a")
	b()
}

func main() {
	a()
}
