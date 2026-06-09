package main

import (
	"fmt"
	"unsafe"
)

func main() {
	// Дескриптор слайса: 3 поля, 24 байта на 64-bit
	s := make([]int, 3, 6)
	fmt.Println("sizeof(slice header):", unsafe.Sizeof(s), "bytes")
	fmt.Println("len:", len(s), "cap:", cap(s))

	// Data указывает на backing array
	fmt.Printf("Data ptr: %p\n", unsafe.SliceData(s))
	fmt.Printf("&s[0]:    %p\n", &s[0])

	// make([]T, len, cap): len доступных элементов, cap — всего в backing array
	s2 := make([]int, 3, 6)
	s2[0], s2[1], s2[2] = 10, 20, 30
	fmt.Println("s2:", s2, "len:", len(s2), "cap:", cap(s2))

	// append использует «скрытую» capacity без реаллокации
	s2 = append(s2, 40)
	fmt.Println("после append(40):", s2, "len:", len(s2), "cap:", cap(s2))

	// При реаллокации Data меняется
	before := unsafe.SliceData(s2)
	for i := 0; i < 10; i++ {
		s2 = append(s2, i)
	}
	after := unsafe.SliceData(s2)
	fmt.Println("ptr до реаллокации:", before)
	fmt.Println("ptr после реаллокации:", after)
	fmt.Println("ptr изменился:", before != after)

	// Передача слайса в функцию копирует дескриптор, не данные
	original := []int{1, 2, 3, 4, 5}
	modifyElements(original)
	fmt.Println("original после modifyElements:", original)

	appendInside(original)
	fmt.Println("original после appendInside:", original, "len:", len(original))
}

func modifyElements(s []int) {
	for i := range s {
		s[i] *= 2
	}
}

func appendInside(s []int) {
	s = append(s, 999)
	fmt.Println("внутри appendInside, len:", len(s))
}
