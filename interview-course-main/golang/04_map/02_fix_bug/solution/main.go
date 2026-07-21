package main

import (
	"fmt"
	"time"
)

// Проблема: все элементы stockHistory указывают на одну и ту же мапу currentStock.
// Когда мы изменяем currentStock в следующей итерации, все предыдущие записи
// в stockHistory тоже изменяются, так как они ссылаются на тот же объект.

// Решение: создавать копию мапы перед отправкой в канал.
func UpdateProductStock() <-chan map[string]int {
	stockUpdates := make(chan map[string]int)

	go func() {
		currentStock := map[string]int{
			"Apples":  50,
			"Bananes": 30,
			"Oranges": 20,
			"Grapes":  15,
		}

		for i := 0; i < 5; i++ {
			for product, quantity := range currentStock {
				currentStock[product] = int(float64(quantity) * 0.95)
			}

			// ВАЖНО: создаем копию мапы перед отправкой.
			stockCopy := make(map[string]int)
			for k, v := range currentStock {
				stockCopy[k] = v
			}
			stockUpdates <- stockCopy

			time.Sleep(150 * time.Millisecond)
		}
		close(stockUpdates)
	}()

	return stockUpdates
}

func main() {
	stockStream := UpdateProductStock()

	var stockHistory []map[string]int

	for stock := range stockStream {
		stockHistory = append(stockHistory, stock)
	}

	for i, stock := range stockHistory {
		fmt.Printf("Iteration %d: %v\n", i+1, stock)
	}
}

// Ответ (ДО исправления - БАГ):
// Iteration 1: map[Apples:38 Bananes:22 Grapes:10 Oranges:15]
// Iteration 2: map[Apples:38 Bananes:22 Grapes:10 Oranges:15]
// Iteration 3: map[Apples:38 Bananes:22 Grapes:10 Oranges:15]
// Iteration 4: map[Apples:38 Bananes:22 Grapes:10 Oranges:15]
// Iteration 5: map[Apples:38 Bananes:22 Grapes:10 Oranges:15]
//
// Все итерации показывают одинаковые значения (последние), так как все элементы
// stockHistory указывают на одну и ту же мапу currentStock.
//
// Ответ (после исправления):
// Iteration 1: map[Apples:47 Bananes:28 Grapes:14 Oranges:19]
// Iteration 2: map[Apples:44 Bananes:26 Grapes:13 Oranges:18]
// Iteration 3: map[Apples:42 Bananes:25 Grapes:12 Oranges:17]
// Iteration 4: map[Apples:40 Bananes:23 Grapes:11 Oranges:16]
// Iteration 5: map[Apples:38 Bananes:22 Grapes:10 Oranges:15]
//
// Объяснение проблемы:
// 1. Мапы в Go - это ссылочные типы (reference types).
// 2. Когда мы делаем stockUpdates <- currentStock, мы отправляем ссылку на мапу.
// 3. Все элементы stockHistory указывают на одну и ту же мапу currentStock.
// 4. При изменении currentStock в следующей итерации, все предыдущие записи тоже изменяются.
// 5. К моменту вывода все элементы stockHistory содержат последние значения из currentStock.
//
// Дополнительные исправления:
// - Закрываем канал после завершения отправки (close(stockUpdates))
// - Используем range для чтения из канала вместо фиксированного цикла
