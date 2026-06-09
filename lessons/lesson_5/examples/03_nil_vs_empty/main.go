package main

import (
	"encoding/json"
	"fmt"
	"unsafe"
)

func main() {
	// Четыре способа получить «пустой или nil» слайс
	var nilSlice []string          // nil slice
	explicitNil := []string(nil)   // nil slice
	emptyLit := []string{}         // empty slice
	emptyMake := make([]string, 0) //empty slice

	if len(explicitNil) == 0 {
		fmt.Println()
	}

	printSliceInfo("var nilSlice []string", nilSlice)
	printSliceInfo("[]string(nil)", explicitNil)
	printSliceInfo("[]string{}", emptyLit)
	printSliceInfo("make([]string, 0)", emptyMake)

	// append и range работают одинаково для nil и empty
	nilSlice = append(nilSlice, "hello")
	fmt.Println("append к nil-слайсу:", nilSlice)

	var iterNil []int
	count := 0
	for range iterNil {
		count++
	}
	fmt.Println("range по nil-слайсу, итераций:", count)

	// JSON-сериализация: ключевое различие
	var ns []int
	es := []int{}
	jsonNil, _ := json.Marshal(ns)
	jsonEmpty, _ := json.Marshal(es)
	fmt.Println("JSON nil-слайса:", string(jsonNil))
	fmt.Println("JSON empty-слайса:", string(jsonEmpty))
}

func printSliceInfo(label string, s []string) {
	fmt.Printf("%-35s → nil=%-5t empty=%-5t len=%-2d ptr=%v\n",
		label,
		s == nil,
		len(s) == 0,
		len(s),
		unsafe.SliceData(s),
	)
}
