/*
Домашнее задание 3: Возведение в степень

Реализуй функцию power(base, exp int) int — целочисленное возведение в степень.

Требования:
  - exp >= 0 (для отрицательного exp — паника или возврат ошибки, на выбор)
  - power(x, 0) = 1 для любого x (включая 0)
  - Реализуй наивную версию: умножай base в цикле exp раз

Проверь на:
  power(2, 0)  = 1
  power(2, 10) = 1024
  power(3, 5)  = 243
  power(10, 6) = 1000000
  power(0, 5)  = 0
  power(1, 100) = 1

Усложнение — быстрое возведение в степень (бинарный метод):
  Идея: base^exp = (base^(exp/2))^2 если exp чётный
              = base * base^(exp-1) если exp нечётный
  Это O(log n) вместо O(n) умножений.

  Реализуй powerFast(base, exp int) int и сравни количество шагов
  (добавь счётчик умножений в обе функции).
*/

package main

import "fmt"

func power(base, exp int) int {
	if exp < 0 {
		panic("отрицательная степень не поддерживается")
	}
	// TODO: наивная версия — цикл exp раз
	return 0
}

func powerFast(base, exp int) int {
	if exp < 0 {
		panic("отрицательная степень не поддерживается")
	}
	// TODO: быстрое возведение в степень (бинарный метод)
	return 0
}

func main() {
	tests := [][2]int{{2, 0}, {2, 10}, {3, 5}, {10, 6}, {0, 5}, {1, 100}}
	for _, t := range tests {
		base, exp := t[0], t[1]
		slow := power(base, exp)
		fast := powerFast(base, exp)
		match := "✓"
		if slow != fast {
			match = "✗ РАСХОЖДЕНИЕ"
		}
		fmt.Printf("power(%d, %d) = %d  %s\n", base, exp, slow, match)
	}
}
