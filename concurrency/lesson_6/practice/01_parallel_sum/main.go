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

	//  len(data) / parts = 1000 /  4 = 250
	// 0:250, 250:500, 500:750, 750:1000
	// 0*250, i+1 * 250
	r := 0
	resSlice := make([]int, parts)
	wg := &sync.WaitGroup{}
	wg.Add(parts)
	for i := 0; i < parts; i++ {
		l := len(data) / parts
		go func() {
			defer wg.Done()
			for _, v := range data[i*l : (i+1)*l] {
				resSlice[i] += v
			}
			fmt.Println(resSlice[i])
		}()
	}
	wg.Wait()
	for _, v := range resSlice {
		r += v
	}
	return r
}

func main() {
	// Засекаем время выполнения
	start := time.Now()

	data := make([]int, 2000000000)
	for i := range data {
		data[i] = i + 1
	}
	fmt.Println("sum:", parallelSum(data, 16)) // ожидается 500500

	// Выводим время выполнения
	elapsed := time.Since(start)
	fmt.Printf("Время выполнения: %v\n", elapsed)
}
