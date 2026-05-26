package main

import "fmt"

func main() {
	// for с одним условием — аналог while в других языках
	n := 1
	for n < 100 {
		n *= 2
	}
	fmt.Println("Первая степень двойки >= 100:", n)

	// Суммирование цифр числа
	num := 12345
	sum := 0
	for num > 0 {
		sum += num % 10
		num /= 10
	}
	fmt.Println("Сумма цифр:", sum) // 15
}
