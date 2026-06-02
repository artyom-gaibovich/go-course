package main

import (
	"fmt"
)

func accessToElement1() {
	data := make([]int, 3)
	fmt.Println(data[4]) // panic
}

func accessToElement2() {
	data := make([]int, 3, 6)
	fmt.Println(data[4]) // panic
}

func accessToNilSlice1() {
	var data []int
	_ = data[0] // panic
}

func accessToNilSlice2() {
	var data []int
	data[0] = 10 // panic
}

func appendToNilSlice() {
	var data []int
	data = append(data, 10)
}

func rangeByNilSlice() {
	var data []int
	for range data {
	}
}

func makeZeroSlice() {
	data := make([]int, 0)
	fmt.Println(len(data)) // 0
	fmt.Println(cap(data)) // 0
}

func main() {
	//var s []int  // len=0, cap=0, ptr=nil
	//b := []int{} // len=0, cap=0, ptr!=nil
	//
	//a := make([]int, 5) // len=5, cap=5, элементы = 0
	//
	//d := make([]int, 3, 6) // len=3, cap=6, элементы = 0
	//
	//data := make([]int, 0)
	//fmt.Println(len(data)) // 0
	//fmt.Println(cap(data)) // 0
	//
	//data = append(data, 1)
	//
	//fmt.Println(s, b, a, d)
	//
	//// len = capacity
	//
	//// make([]T, 0, n)

	s := []int{1, 2, 3}
	fill(s)
	fmt.Println(s)
}

func fill(s []int) {
	s[0] = 99
	s = append(s, 1)

}

//func makeSlice() {
//	_ = make([]int, -5)    // compilation error
//	_ = make([]int, 10, 5) // compilation error
//
//	size := -5
//	_ = make([]int, size) // panic
//
//	size = 5
//	_ = make([]int, size*2, size) // panic
//}
//
//func accessToElement3() {
//	data := make([]int, 3, 6)
//	_ = data[-1] // compilation error
//}
