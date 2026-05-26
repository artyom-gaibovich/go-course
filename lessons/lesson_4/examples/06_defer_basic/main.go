package main

import "fmt"

func main() {
	// defer откладывает вызов функции до возврата из текущей функции.
	// Аргументы вычисляются СРАЗУ при объявлении defer, не при выполнении.

	fmt.Println("начало")

	defer fmt.Println("defer 1")
	defer fmt.Println("defer 2")
	defer fmt.Println("defer 3")

	fmt.Println("конец")
	// Вывод:
	// начало
	// конец
	// defer 3   ← LIFO: последний зарегистрированный выполняется первым
	// defer 2
	// defer 1
}
