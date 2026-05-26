package main

import "fmt"

func main() {
	score := 74

	// switch без выражения работает как цепочка if-else
	switch {
	case score >= 90:
		fmt.Println("Отлично")
	case score >= 75:
		fmt.Println("Хорошо")
	case score >= 60:
		fmt.Println("Удовлетворительно")
	default:
		fmt.Println("Неудовлетворительно")
	}

	// switch с инициализатором
	switch x := score * 2; {
	case x > 150:
		fmt.Println("Двойной балл высокий:", x)
	default:
		fmt.Println("Двойной балл:", x)
	}
}
