package main

import "fmt"

func main() {
	fmt.Println("=== Опасный append: разделяемый backing array ===")

	// data и sub шарят один backing array, у data есть свободная capacity
	data := make([]int, 4, 6)
	data[0], data[1], data[2], data[3] = 1, 2, 3, 4
	sub := data[:2]

	fmt.Println("до append:")
	fmt.Println("  data:", data, "len:", len(data), "cap:", cap(data))
	fmt.Println("  sub:", sub, "len:", len(sub), "cap:", cap(sub))

	// append к sub использует свободную capacity backing array
	// и ПЕРЕТИРАЕТ data[2] и data[3]
	sub = append(sub, 100, 200)
	fmt.Println("после sub = append(sub, 100, 200):")
	fmt.Println("  data:", data, "← data[2] и data[3] изменились!")
	fmt.Println("  sub:", sub)

	fmt.Println()
	fmt.Println("=== Защита #1: full slice expression data[:2:2] ===")

	data2 := make([]int, 4, 6)
	data2[0], data2[1], data2[2], data2[3] = 1, 2, 3, 4
	safe := data2[:2:2]

	safe = append(safe, 100, 200)
	fmt.Println("data2 после append к safe [:2:2]:", data2, "← не изменился!")
	fmt.Println("safe:", safe)

	fmt.Println()
	fmt.Println("=== Защита #2: явный copy ===")

	data3 := make([]int, 4, 6)
	data3[0], data3[1], data3[2], data3[3] = 1, 2, 3, 4
	copied := make([]int, 2)
	copy(copied, data3)

	copied = append(copied, 100, 200)
	fmt.Println("data3 после append к copied:", data3, "← не изменился!")
	fmt.Println("copied:", copied)

	fmt.Println()
	fmt.Println("=== Защита #3: идиома append([]T(nil), s...) ===")

	data4 := []int{1, 2, 3, 4}
	clone := append([]int(nil), data4[:4]...)
	clone = append(clone, 100, 200)
	fmt.Println("data4 после append к clone:", data4, "← не изменился!")
	fmt.Println("clone:", clone)
}
