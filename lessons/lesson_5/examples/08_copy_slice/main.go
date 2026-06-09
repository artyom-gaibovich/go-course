package main

import "fmt"

func main() {
	fmt.Println("=== copy: базовая семантика ===")

	src := []int{1, 2, 3, 4, 5}
	dst := make([]int, len(src))
	n := copy(dst, src)
	fmt.Println("скопировано элементов:", n)
	fmt.Println("dst:", dst)

	// copy ограничен min(len(dst), len(src))
	small := make([]int, 3)
	n2 := copy(small, src)
	fmt.Println("copy в маленький dst:", n2, "→", small)

	large := make([]int, 8)
	n3 := copy(large, src)
	fmt.Println("copy в большой dst:", n3, "→", large)

	fmt.Println()
	fmt.Println("=== Частая ошибка: copy в nil/пустой dst ===")

	var nilDst []int
	n4 := copy(nilDst, src)
	fmt.Println("copy в nil dst:", n4, "элементов (ноль!)")

	emptyDst := make([]int, 0, 10)
	n5 := copy(emptyDst, src)
	fmt.Println("copy в пустой dst (len=0):", n5, "элементов (ноль!)")

	fmt.Println()
	fmt.Println("=== copy строки в []byte ===")

	str := "Hello, Go!"
	buf := make([]byte, len(str))
	copy(buf, str)
	fmt.Println("строка → []byte:", string(buf))

	fmt.Println()
	fmt.Println("=== copy при перекрывающихся диапазонах ===")

	data := []int{1, 2, 3, 4, 5}
	copy(data[1:], data[:4])
	fmt.Println("copy(data[1:], data[:4]):", data)

	data2 := []int{1, 2, 3, 4, 5}
	copy(data2[:4], data2[1:])
	fmt.Println("copy(data2[:4], data2[1:]):", data2)

	fmt.Println()
	fmt.Println("=== Удаление элемента из середины через copy ===")

	nums := []int{10, 20, 30, 40, 50}
	i := 2
	nums = append(nums[:i], nums[i+1:]...)
	fmt.Println("после удаления index=2:", nums)
}
