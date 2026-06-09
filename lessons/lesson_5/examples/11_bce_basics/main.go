package main

import "fmt"

// Запустите с флагом для наблюдения за BCE:
// go run -gcflags="-d=ssa/check_bce" re.go

// BCE работает: компилятор знает i < len(s) из условия range
func sumRange(s []int) int {
	total := 0
	for i := range s {
		total += s[i]
	}
	return total
}

// BCE работает: i < len(s) в условии цикла
func sumFor(s []int) int {
	total := 0
	for i := 0; i < len(s); i++ {
		total += s[i]
	}
	return total
}

// BCE НЕ работает: компилятор не знает, что i < len(s)
// (size может быть больше len(values))
func sumWithSize(values []int, size int) int {
	total := 0
	for i := 0; i < size; i++ {
		total += values[i]
	}
	return total
}

// BCE с хинтом: одно обращение до цикла доказывает компилятору bounds
func sumWithSizeOptimized(values []int, size int) int {
	if size == 0 {
		return 0
	}
	_ = values[size-1]
	total := 0
	for i := 0; i < size; i++ {
		total += values[i]
	}
	return total
}

// BCE работает для обоих слайсов через min(len)
func dotProduct(a, b []int) int {
	n := min(len(a), len(b))
	result := 0
	for i := range n {
		result += a[i] * b[i]
	}
	return result
}

func main() {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	fmt.Println("sumRange:", sumRange(data))
	fmt.Println("sumFor:", sumFor(data))
	fmt.Println("sumWithSize:", sumWithSize(data, 10))
	fmt.Println("sumWithSizeOptimized:", sumWithSizeOptimized(data, 10))

	a := []int{1, 2, 3}
	b := []int{4, 5, 6}
	fmt.Println("dotProduct:", dotProduct(a, b))
}
