package main

import (
	"fmt"
	"strconv"
)

func main() {
	// Инициализирующий оператор в if: переменная видна только внутри блока

	if n, err := strconv.Atoi("42"); err == nil {
		fmt.Println("Число:", n)
	} else {
		fmt.Println("Ошибка:", err)
	}

	// n и err здесь недоступны — за пределами if они не существуют
	input := "не число"
	if val, err := strconv.Atoi(input); err != nil {
		fmt.Printf("Не удалось распарсить %q: %v\n", input, err)
	} else {
		fmt.Println("Значение:", val)
	}
}
