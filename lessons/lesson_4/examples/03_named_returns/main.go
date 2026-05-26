package main

import "fmt"

// Named returns: возвращаемые переменные объявлены в сигнатуре.
// Они инициализируются нулевыми значениями и доступны в теле функции.
func minMax(nums []int) (min, max int) {
	if len(nums) == 0 {
		return // "голый" return — возвращает min и max как есть (оба 0)
	}
	min, max = nums[0], nums[0]
	for _, n := range nums[1:] {
		if n < min {
			min = n
		}
		if n > max {
			max = n
		}
	}
	return // возвращает текущие значения min и max
}

// Named returns полезны для документирования смысла возвращаемых значений,
// особенно когда типы одинаковые — (min, max int) читается лучше (int, int).
func divide(a, b float64) (result float64, err error) {
	if b == 0 {
		err = fmt.Errorf("деление на ноль")
		return
	}
	result = a / b
	return
}

func main() {
	nums := []int{3, 1, 4, 1, 5, 9, 2, 6}
	min, max := minMax(nums)
	fmt.Printf("min=%d max=%d\n", min, max)

	res, err := divide(10, 4)
	fmt.Printf("10/4 = %.2f, err = %v\n", res, err)
}
