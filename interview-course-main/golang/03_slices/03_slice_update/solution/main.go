package main

import "fmt"

// modifyElement изменяет элемент по индексу.
// Слайс содержит указатель на массив, поэтому изменения видны снаружи.
func modifyElement(slice []int) {
	slice[1] = 999
}

// addElement добавляет элемент и изменяет первый элемент.
// Но append возвращает новый слайс, который мы присваиваем локальной переменной slice.
func addElement(slice []int) {
	slice = append(slice, 100)
	slice[0] = 888
	fmt.Println("Внутри addElement:", slice)
}

func main() {
	original := []int{10, 20, 30}

	fmt.Println("До modifyElement:", original)
	modifyElement(original)
	fmt.Println("После modifyElement:", original)

	fmt.Println("До addElement:", original)
	addElement(original)
	fmt.Println("После addElement:", original)
}

// Ответ:
// До modifyElement: [10 20 30]
// После modifyElement: [10 999 30]
// До addElement: [10 999 30]
// Внутри addElement: [888 999 30 100]
// После addElement: [10 999 30]
//
// Объяснение:
// 1. modifyElement изменяет элемент по индексу через указатель на массив.
//    Изменения видны в original, так как слайс содержит указатель на тот же массив.
//
// 2. addElement делает append, который возвращает новый слайс.
//    Мы присваиваем его локальной переменной slice, но не возвращаем из функции.
//    original не изменяется, так как мы не модифицировали его напрямую.
//
// 3. slice[0] = 888 изменяет первый элемент в новом слайсе (который может быть
//    в том же массиве, если capacity достаточен, или в новом массиве).
//    Но это не влияет на original, так как мы работаем с локальной копией слайса.
//
// Чтобы изменения были видны снаружи, нужно вернуть новый слайс:
// func addElement(slice []int) []int {
//     return append(slice, 100)
// }
