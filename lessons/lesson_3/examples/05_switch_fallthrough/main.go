package main

import "fmt"

func main() {
	n := 1

	// По умолчанию Go НЕ проваливается в следующий case (в отличие от C).
	// fallthrough явно передаёт управление следующему case (без проверки условия).
	switch n {
	case 1:
		fmt.Println("один")
		fallthrough
	case 2:
		fmt.Println("два")
		fallthrough
	case 3:
		fmt.Println("три")
	case 4:
		fmt.Println("четыре")
	}
	// Вывод: два, три (начало с case 2, fallthrough в case 3)
}
