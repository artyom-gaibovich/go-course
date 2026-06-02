package main

import "fmt"

func main() {
	// defer откладывает вызов функции до возврата из текущей функции.
	// Аргументы вычисляются СРАЗУ при объявлении defer, не при выполнении.

	v := 10
	fmt.Println("начало")

	defer fmt.Println("defer value", v)
	defer fmt.Println("defer 2")
	defer fmt.Println("defer 3")

	v++
	fmt.Println("first", v)

	// Вывод:
	// начало
	// конец
	// defer 3   ← LIFO: последний зарегистрированный выполняется первым
	// defer 2
	// defer 1
}
