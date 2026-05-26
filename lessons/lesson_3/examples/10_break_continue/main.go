package main

import "fmt"

func main() {
	// continue — пропустить итерацию
	fmt.Println("Нечётные числа до 10:")
	for i := 1; i <= 10; i++ {
		if i%2 == 0 {
			continue
		}
		fmt.Print(i, " ")
	}
	fmt.Println()

	// break — выйти из цикла
	fmt.Println("Числа до первого кратного 6:")
	for i := 1; i <= 20; i++ {
		if i%6 == 0 {
			break
		}
		fmt.Print(i, " ")
	}
	fmt.Println()

	// Метки (labels) для break/continue из вложенного цикла
outer:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if i == 1 && j == 1 {
				break outer // выходим из внешнего цикла
			}
			fmt.Printf("(%d,%d) ", i, j)
		}
	}
	fmt.Println()
}
