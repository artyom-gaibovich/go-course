package main

import "fmt"

func main() {
	// Бесконечный цикл — for без условия
	// Используется когда условие выхода внутри логики цикла
	counter := 0
	for {
		counter++
		if counter == 5 {
			fmt.Println("Достигли 5, выходим")
			break
		}
	}

	// Поиск первого числа, делящегося на 7 и на 11
	n := 1
	for {
		if n%7 == 0 && n%11 == 0 {
			fmt.Println("Первое кратное 7 и 11:", n) // 77
			break
		}
		n++
	}
}
