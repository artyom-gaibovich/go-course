package main

import "fmt"

func main() {
	fruits := []string{"яблоко", "банан", "вишня"}

	// range возвращает индекс и значение
	for i, v := range fruits {
		fmt.Printf("%d: %s\n", i, v)
	}

	// Если индекс не нужен — используем _
	for _, v := range fruits {
		fmt.Println(v)
	}

	// range по строке — итерирует по рунам (Unicode), не байтам
	for i, r := range "Go!" {
		fmt.Printf("индекс=%d руна=%c код=%d\n", i, r, r)
	}

	// range по map (порядок случайный)
	scores := map[string]int{"Алиса": 95, "Боб": 87}
	for name, score := range scores {
		fmt.Printf("%s: %d\n", name, score)
	}
}
