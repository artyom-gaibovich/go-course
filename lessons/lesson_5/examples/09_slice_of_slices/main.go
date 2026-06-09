package main

import (
	"fmt"
	"runtime"
)

type Block struct {
	data []byte
}

func printAllocs(label string) {
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("%-40s %d KB\n", label+":", m.Alloc/1024)
}

func main() {
	fmt.Println("=== Утечка через слайс слайсов ===")
	printAllocs("старт")

	blocks := makeBlocks(100)
	printAllocs("после makeBlocks(100)")

	leaked := blocks[:2]
	blocks = nil
	printAllocs("после blocks=nil, leaked=blocks[:2]")

	_ = leaked

	fmt.Println()
	fmt.Println("=== Правильная очистка через clear ===")
	printAllocs("старт 2")

	blocks2 := makeBlocks(100)
	printAllocs("после makeBlocks(100)")

	clear(blocks2[2:])
	freed := blocks2[:2]
	blocks2 = nil
	printAllocs("после clear + blocks2=nil")

	_ = freed

	fmt.Println()
	fmt.Println("=== Создание матрицы (слайс слайсов) ===")
	matrix := makeMatrix(4, 4)
	for i, row := range matrix {
		for j := range row {
			matrix[i][j] = i*4 + j + 1
		}
	}
	printMatrix(matrix)

	// Подматрица шарит backing array строк
	sub := matrix[1:3]
	sub[0][0] = 99
	fmt.Println("matrix[1][0] после sub[0][0]=99:", matrix[1][0])
}

func makeBlocks(n int) []Block {
	blocks := make([]Block, n)
	for i := range blocks {
		blocks[i] = Block{data: make([]byte, 4096)}
	}
	return blocks
}

func makeMatrix(rows, cols int) [][]int {
	matrix := make([][]int, rows)
	for i := range matrix {
		matrix[i] = make([]int, cols)
	}
	return matrix
}

func printMatrix(m [][]int) {
	for _, row := range m {
		fmt.Println(row)
	}
}
