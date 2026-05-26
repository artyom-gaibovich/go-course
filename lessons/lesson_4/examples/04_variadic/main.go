package main

import "fmt"

// Variadic: функция принимает произвольное количество аргументов одного типа.
// Внутри nums — обычный []int.
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// Variadic параметр всегда последний
func logMessage(level string, parts ...string) {
	fmt.Printf("[%s]", level)
	for _, p := range parts {
		fmt.Print(" ", p)
	}
	fmt.Println()
}

func main() {
	fmt.Println(sum())           // 0
	fmt.Println(sum(1))          // 1
	fmt.Println(sum(1, 2, 3))    // 6
	fmt.Println(sum(1, 2, 3, 4)) // 10

	// Слайс можно передать в variadic с помощью ... (spread)
	nums := []int{10, 20, 30}
	fmt.Println(sum(nums...)) // 60

	logMessage("INFO", "сервер", "запущен")
	logMessage("ERROR", "файл", "не", "найден")
}
