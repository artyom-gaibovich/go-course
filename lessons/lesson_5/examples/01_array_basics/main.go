package main

import "fmt"

func main() {
	// Размер — часть типа: [3]int и [4]int несовместимы
	var a3 [3]int
	var a4 [4]int
	fmt.Println("a3:", a3, "len:", len(a3), "cap:", cap(a3))
	fmt.Println("a4:", a4, "len:", len(a4), "cap:", cap(a4))

	// Способы инициализации
	lit := [5]int{10, 20, 30}
	fmt.Println("lit:", lit)

	inferred := [...]int{100, 200, 300}
	fmt.Println("inferred:", inferred, "len:", len(inferred))

	// Адрес массива == адрес первого элемента
	fmt.Println("&inferred:", &inferred, "&inferred[0]:", &inferred[0])

	// Массивы можно сравнивать поэлементно через ==
	x := [3]int{1, 2, 3}
	y := [3]int{1, 2, 3}
	z := [3]int{1, 2, 4}
	fmt.Println("x==y:", x == y, "x==z:", x == z)

	// Итерация: range делает копию массива
	arr := [3]int{1, 2, 3}
	for i, v := range arr {
		v = v * 100
		_ = v
		_ = i
	}
	fmt.Println("arr после range (v не меняет оригинал):", arr)

	// Правильное изменение через индекс
	for i := range arr {
		arr[i] *= 10
	}
	fmt.Println("arr после arr[i]*=10:", arr)

	// Передача по указателю избегает копирования
	bigArr := [1000]int{}
	bigArr[0] = 42
	modifyByPointer(&bigArr)
	fmt.Println("bigArr[0] после modifyByPointer:", bigArr[0])
}

// 1000 МБ

// modifyByPointer() as value
// 2000 мб
func modifyByPointer(arr *[1000]int) {
	arr[0] = 99
}
