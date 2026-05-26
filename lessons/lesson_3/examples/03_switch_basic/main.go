package main

import "fmt"

func main() {
	day := 3

	switch day {
	case 1:
		fmt.Println("Понедельник")
	case 2:
		fmt.Println("Вторник")
	case 3:
		fmt.Println("Среда")
	case 4:
		fmt.Println("Четверг")
	case 5:
		fmt.Println("Пятница")
	case 6, 7:
		fmt.Println("Выходной") // несколько значений в одном case
	default:
		fmt.Println("Неверный день")
	}

	// switch со строкой
	lang := "Go"
	switch lang {
	case "Go":
		fmt.Println("Быстрый и лаконичный")
	case "Python":
		fmt.Println("Динамичный и читаемый")
	default:
		fmt.Println("Хороший выбор!")
	}
}
