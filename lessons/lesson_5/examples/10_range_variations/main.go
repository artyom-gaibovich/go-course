package main

import "fmt"

func main() {
	fmt.Println("=== range по массиву делает копию ===")
	arr := [3]int{10, 20, 30}
	for i, v := range arr {
		v *= 100
		_ = v
		_ = i
	}
	fmt.Println("arr после range (v не менял оригинал):", arr)

	fmt.Println()
	fmt.Println("=== range по &array — без копии ===")
	arr2 := [3]int{10, 20, 30}
	for i := range &arr2 {
		arr2[i] *= 10
	}
	fmt.Println("arr2 после range &arr2:", arr2)

	fmt.Println()
	fmt.Println("=== range по arr[:] — без копии ===")
	arr3 := [3]int{1, 2, 3}
	slice := arr3[:]
	for i := range slice {
		slice[i] *= 5
	}
	fmt.Println("arr3 после изменений через slice:", arr3)

	fmt.Println()
	fmt.Println("=== range фиксирует len в момент старта ===")
	growing := []int{1, 2, 3}
	iterations := 0
	for range growing {
		growing = append(growing, 99)
		iterations++
	}
	fmt.Println("итераций:", iterations, "(зафиксировано при старте: 3)")
	fmt.Println("growing:", growing)

	fmt.Println()
	fmt.Println("=== nil-указатель на массив: len/cap работают, разыменование — нет ===")
	var p *[4]int
	fmt.Println("len(nil ptr):", len(p), "cap(nil ptr):", cap(p))
	for idx := range p {
		fmt.Println("for range p, idx:", idx)
	}

	fmt.Println()
	fmt.Println("=== Go 1.22: каждая итерация — новая переменная v ===")
	funcs := make([]func(), 3)
	values := [3]int{10, 20, 30}
	for i, v := range values {
		v := v
		funcs[i] = func() { fmt.Print(v, " ") }
	}
	for _, f := range funcs {
		f()
	}
	fmt.Println("← каждая closure захватила свою копию v")
}
