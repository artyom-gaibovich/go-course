package main

import (
	"fmt"
	"unsafe"
)

func main() {
	// Наблюдаем рост capacity при последовательных append
	s := make([]int, 0)
	prevCap := 0
	for i := 0; i < 65536; i++ {
		s = append(s, i)
		if cap(s) != prevCap {
			fmt.Printf("len=%-3d cap=%-4d  (рост с %d)\n", len(s), cap(s), prevCap)
			prevCap = cap(s)
		}
	}

	fmt.Println()

	// append без реаллокации: len < cap
	withRoom := make([]int, 2, 5)
	withRoom[0], withRoom[1] = 10, 20
	ptrBefore := unsafe.SliceData(withRoom)
	withRoom = append(withRoom, 30)
	ptrAfter := unsafe.SliceData(withRoom)
	fmt.Println("append без реаллокации:")
	fmt.Println("  ptr изменился:", ptrBefore != ptrAfter, "len:", len(withRoom), "cap:", cap(withRoom))

	// append с реаллокацией: len == cap
	full := []int{1, 2, 3}
	ptrBefore2 := unsafe.SliceData(full)
	full = append(full, 4)
	ptrAfter2 := unsafe.SliceData(full)
	fmt.Println("append с реаллокацией:")
	fmt.Println("  ptr изменился:", ptrBefore2 != ptrAfter2, "len:", len(full), "cap:", cap(full))

	// Предвыделение: 0 реаллокаций
	preallocated := make([]int, 0, 1000)
	ptrStart := unsafe.SliceData(preallocated)
	for i := 0; i < 1000; i++ {
		preallocated = append(preallocated, i)
	}
	ptrEnd := unsafe.SliceData(preallocated)
	fmt.Println("1000 append с предвыделением, ptr изменился:", ptrStart != ptrEnd)

	// append нескольких элементов сразу
	a := []int{1, 2}
	b := []int{10, 20, 30}
	a = append(a, b...)
	fmt.Println("append нескольких элементов:", a)
}
