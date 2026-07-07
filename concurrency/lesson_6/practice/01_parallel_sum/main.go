package main

import (
	"fmt"
	"sync"
	"time"
)

func parallelSum(data []int, parts int) int {
	// TODO: разбить data на parts частей
	// TODO: запустить горутину на каждую часть, результат в results[i]
	// TODO: дождаться всех через wg.Wait() и сложить results

}

func main() {
	// Засекаем время выполнения
	start := time.Now()

	data := make([]int, 64)
	for i := range data {
		data[i] = i + 1
	}
	fmt.Println("sum:", parallelSum(data, 4)) // ожидается 500500

	// Выводим время выполнения
	elapsed := time.Since(start)
	fmt.Printf("Время выполнения: %v\n", elapsed)
}
