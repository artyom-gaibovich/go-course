package main

import "fmt"

func main() {
	age := 20

	if age < 18 {
		fmt.Println("несовершеннолетний")
	} else if age < 30 {
		fmt.Println("взрослый 1")
	} else if age < 65 {
		fmt.Println("взрослый 2")
	} else {
		fmt.Println("пенсионер")
	}

	// В Go скобки вокруг условия не нужны, фигурные — обязательны
	x := 10
	if x > 0 {
		fmt.Println("положительное")
	}
}
