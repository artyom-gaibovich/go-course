// Задача: Что выведет программа и почему?

package main

import "fmt"

func main() {
	data := []int{10, 20, 30, 40}
	fmt.Println("Изначальный слайс:", data)

	modify(data[:2])
	fmt.Println("Слайс после изменений:", data)
}

func modify(slice []int) {
	slice = append(slice, 50, 60)
	fmt.Println("Слайс в функции модификации:", slice)
}
