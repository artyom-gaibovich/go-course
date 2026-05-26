package main

import "fmt"

// Базовая функция: имя, параметры, тип возврата
func add(a, b int) int {
	return a + b
}

// Параметры одного типа можно сгруппировать: (a, b int) вместо (a int, b int)
func multiply(a, b int) int {
	return a * b
}

// Функция без возвращаемого значения
func greet(name string) {
	fmt.Println("Привет,", name)
}

// Функция с несколькими параметрами разных типов
func describe(name string, age int, active bool) string {
	status := "неактивен"
	if active {
		status = "активен"
	}
	return fmt.Sprintf("%s (%d лет) — %s", name, age, status)
}

func main() {
	fmt.Println(add(3, 4))
	fmt.Println(multiply(6, 7))
	greet("Мир")
	fmt.Println(describe("Алиса", 30, true))
}
