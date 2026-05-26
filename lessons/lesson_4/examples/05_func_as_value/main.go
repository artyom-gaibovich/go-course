package main

import "fmt"

// Функции в Go — first-class citizens: их можно присваивать, передавать, возвращать.

// Тип функции: func(int) int
func apply(nums []int, f func(int) int) []int {
	result := make([]int, len(nums))
	for i, n := range nums {
		result[i] = f(n)
	}
	return result
}

// Функция возвращает функцию — замыкание (closure)
func multiplier(factor int) func(int) int {
	return func(n int) int {
		return n * factor // factor "захватывается" из внешней области видимости
	}
}

func main() {
	// Анонимная функция
	double := func(n int) int { return n * 2 }
	fmt.Println(double(5)) // 10

	nums := []int{1, 2, 3, 4, 5}
	fmt.Println(apply(nums, double))                        // [2 4 6 8 10]
	fmt.Println(apply(nums, func(n int) int { return n * n })) // [1 4 9 16 25]

	// Замыкание: каждый вызов multiplier создаёт свою копию factor
	triple := multiplier(3)
	tenX := multiplier(10)
	fmt.Println(triple(7))  // 21
	fmt.Println(tenX(7))    // 70

	// Немедленно вызываемая функция (IIFE)
	result := func(a, b int) int {
		return a + b
	}(3, 4)
	fmt.Println(result) // 7
}
