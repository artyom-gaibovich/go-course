/*
Домашнее задание 2: Наибольший общий делитель (НОД)

Реализуй функцию gcd(a, b int) int — наибольший общий делитель двух чисел.

Используй алгоритм Евклида:
  НОД(a, 0) = a
  НОД(a, b) = НОД(b, a mod b)

Требования:
  - Работает с положительными числами
  - gcd(a, 0) и gcd(0, b) → возвращай ненулевое число
  - Реализуй рекурсивную версию

Проверь на:
  gcd(48, 18)   = 6
  gcd(100, 75)  = 25
  gcd(17, 13)   = 1  (взаимно простые)
  gcd(0, 5)     = 5
  gcd(12, 12)   = 12

Усложнение: реализуй lcm(a, b int) int — наименьшее общее кратное.
  Формула: lcm(a, b) = a * b / gcd(a, b)
  lcm(4, 6)   = 12
  lcm(21, 14) = 42
*/

package main

import "fmt"

func gcd(a, b int) int {
	// TODO: алгоритм Евклида
	return 0
}

func lcm(a, b int) int {
	// TODO: используй gcd
	return 0
}

func main() {
	gcdTests := [][2]int{{48, 18}, {100, 75}, {17, 13}, {0, 5}, {12, 12}}
	for _, t := range gcdTests {
		fmt.Printf("gcd(%d, %d) = %d\n", t[0], t[1], gcd(t[0], t[1]))
	}

	fmt.Println()
	lcmTests := [][2]int{{4, 6}, {21, 14}, {3, 5}}
	for _, t := range lcmTests {
		fmt.Printf("lcm(%d, %d) = %d\n", t[0], t[1], lcm(t[0], t[1]))
	}
}
