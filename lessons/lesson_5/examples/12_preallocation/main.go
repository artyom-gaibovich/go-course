package main

import (
	"fmt"
	"runtime"
)

func allocsOf(fn func()) uint64 {
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	fn()
	runtime.ReadMemStats(&after)
	return after.Mallocs - before.Mallocs
}

// Без предвыделения: log2(N) реаллокаций
func filterWithout(input []int) []int {
	result := []int{}
	for _, v := range input {
		if v%2 == 0 {
			result = append(result, v)
		}
	}
	return result
}

// С предвыделением: максимум 1 аллокация
func filterWith(input []int) []int {
	result := make([]int, 0, len(input))
	for _, v := range input {
		if v%2 == 0 {
			result = append(result, v)
		}
	}
	return result
}

func main() {
	const N = 10000
	input := make([]int, N)
	for i := range input {
		input[i] = i
	}

	a1 := allocsOf(func() { filterWithout(input) })
	a2 := allocsOf(func() { filterWith(input) })
	fmt.Printf("filterWithout: ~%d аллокаций\n", a1)
	fmt.Printf("filterWith:    ~%d аллокаций\n", a2)

	fmt.Println()
	fmt.Println("=== clear: memclr-оптимизация ===")
	buf := make([]int, 10)
	for i := range buf {
		buf[i] = i + 1
	}
	fmt.Println("до clear:", buf)
	clear(buf)
	fmt.Println("после clear:", buf)

	fmt.Println()
	fmt.Println("=== Рост capacity: наблюдение за реаллокациями ===")
	s := []int{}
	reallocs := 0
	prevCap := 0
	for i := 0; i < 1024; i++ {
		s = append(s, i)
		if cap(s) != prevCap {
			reallocs++
			prevCap = cap(s)
		}
	}
	fmt.Printf("1024 append без предвыделения: %d реаллокаций\n", reallocs)

	s2 := make([]int, 0, 1024)
	reallocs2 := 0
	prevCap2 := 0
	for i := 0; i < 1024; i++ {
		s2 = append(s2, i)
		if cap(s2) != prevCap2 {
			reallocs2++
			prevCap2 = cap(s2)
		}
	}
	fmt.Printf("1024 append с предвыделением: %d реаллокаций\n", reallocs2)

	fmt.Println()
	fmt.Println("=== Конвертация слайса в массив (Go 1.20+) ===")
	sl := []int{10, 20, 30, 40, 50}
	arr := [3]int(sl[:3])
	sl[0] = 99
	fmt.Println("sl:", sl)
	fmt.Println("arr (копия):", arr, "← sl[0] не влияет на arr")

	arrPtr := (*[3]int)(sl)
	sl[1] = 88
	fmt.Println("arrPtr[1] после sl[1]=88:", arrPtr[1], "← разделяет память")
}
