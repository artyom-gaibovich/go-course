/*
Задание 4: Сумма цифр числа

Напиши функцию sumDigits(n int) int, которая возвращает сумму цифр числа n.
Число может быть отрицательным — знак игнорируй.

Примеры:
  sumDigits(12345) → 15
  sumDigits(0)     → 0
  sumDigits(-987)  → 24

Алгоритм: в цикле забирай последнюю цифру через n % 10, затем n /= 10.

Дополнительно: выведи все промежуточные шаги (какая цифра добавляется на каждом шаге).
*/

package main

import "fmt"

func sumDigits(n int) int {
	if n < 0 {
		n = -n
	}
	sum := 0
	// TODO: реализуй цикл
	return sum
}

func main() {
	tests := []int{12345, 0, -987, 999}
	for _, n := range tests {
		fmt.Printf("sumDigits(%d) = %d\n", n, sumDigits(n))
	}
}
