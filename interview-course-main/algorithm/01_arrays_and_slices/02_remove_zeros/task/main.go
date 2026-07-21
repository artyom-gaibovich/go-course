package main

// Задача: Удаление всех нулей из слайса
//
// Дан слайс целых чисел. Напишите функцию remove, удаляющую все нули.
// Функция должна изменять слайс in-place и возвращать новый слайс нужной длины.
//
// Требования:
// - Изменения in-place (не создавать новый массив)
// - Сохранить порядок ненулевых элементов
// - Оценить временную и пространственную сложность
//
// Примеры:
// remove([]) → []
// remove([0]) → []
// remove([1,0,0,2]) → [1,2]
// remove([0,0,1,2,3]) → [1,2,3]
// remove([1,2,3]) → [1,2,3]

func remove(in []int) []int {
	// TODO: реализуйте функцию
	return nil
}

func main() {
	test1 := []int{}
	println("Test 1:", remove(test1)) // []

	test2 := []int{0}
	println("Test 2:", remove(test2)) // []

	test3 := []int{1, 0, 0, 2}
	println("Test 3:", remove(test3)) // [1, 2]

	test4 := []int{0, 0, 1, 2, 3}
	println("Test 4:", remove(test4)) // [1, 2, 3]

	test5 := []int{1, 2, 3}
	println("Test 5:", remove(test5)) // [1, 2, 3]
}
