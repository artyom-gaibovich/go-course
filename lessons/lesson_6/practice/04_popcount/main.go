/*
Задание 4: Подсчёт установленных битов (Population Count / Hamming Weight)

Реализуй два варианта функции PopCount(n uint64) int:

1. PopCountNaive — наивный: сдвигать n вправо и проверять последний бит в цикле (64 итерации).

2. PopCountKernighan — алгоритм Брайана Кернигана:
   Трюк: n &= n - 1 сбрасывает самый правый единичный бит.
   Цикл выполняется ровно столько раз, сколько единичных битов в числе.

   Пример для n = 0b1010_1100 (4 единицы):
     итерация 1: 0b1010_1100 & 0b1010_1011 = 0b1010_1000
     итерация 2: 0b1010_1000 & 0b1010_0111 = 0b1010_0000
     итерация 3: 0b1010_0000 & 0b1001_1111 = 0b1000_0000
     итерация 4: 0b1000_0000 & 0b0111_1111 = 0b0000_0000 → стоп

Напиши бенчмарк для обоих вариантов в файле main_test.go.
Запусти: go test -bench=. ./...

Проверь:
  PopCount(0)              → 0
  PopCount(1)              → 1
  PopCount(255)            → 8
  PopCount(math.MaxUint64) → 64
*/

package main

import (
	"fmt"
	"math"
)

func PopCountNaive(n uint64) int {
	// TODO: наивный вариант — 64 итерации
	return 0
}

func PopCountKernighan(n uint64) int {
	// TODO: алгоритм Кернигана
	return 0
}

func main() {
	cases := []uint64{0, 1, 7, 255, 0b1010_1100, math.MaxUint64}
	for _, n := range cases {
		naive := PopCountNaive(n)
		kern := PopCountKernighan(n)
		fmt.Printf("PopCount(%064b) = naive:%d kernighan:%d\n", n, naive, kern)
	}
}
