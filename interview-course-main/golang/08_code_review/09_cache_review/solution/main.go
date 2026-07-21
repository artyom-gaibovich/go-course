package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

func LongCalculation(n int) int {
	secondsToSleep := rand.Float64() * float64(n)
	time.Sleep(time.Duration(secondsToSleep))
	return n + 1
}

var (
	cache = make(map[int]int)
	mu    sync.RWMutex
)

// CachedLongCalculation кэширует результаты LongCalculation.
// Исправления:
// 1. Используем один мьютекс для всего кэша (не создаем новый каждый раз).
// 2. Используем RWMutex для оптимизации параллельного чтения.
// 3. Убираем лишний mu.Unlock().
// 4. Добавляем защиту от race condition при параллельных запросах одного ключа.
func CachedLongCalculation(n int) int {
	// Сначала проверяем кэш с RLock (быстрое чтение).
	mu.RLock()
	found, ok := cache[n]
	mu.RUnlock()

	if ok {
		return found
	}

	// Значения нет в кэше, нужно вычислить.
	// Используем Lock для записи.
	mu.Lock()
	// Двойная проверка: возможно, другая горутина уже вычислила значение.
	if found, ok := cache[n]; ok {
		mu.Unlock()
		return found
	}

	// Вычисляем значение (без блокировки, чтобы не блокировать другие горутины).
	mu.Unlock()
	value := LongCalculation(n)

	// Сохраняем результат в кэш.
	mu.Lock()
	// Еще раз проверяем (double-checked locking).
	if found, ok := cache[n]; ok {
		mu.Unlock()
		return found
	}
	cache[n] = value
	mu.Unlock()

	return value
}

func main() {
	nums := []int{5, 10, 22}
	for _, n := range nums {
		val := CachedLongCalculation(n)
		fmt.Printf("LongCalculation(%d) = %d\n", n, val)
	}
}

// Проблемы в исходном коде:
// 1. var mu sync.Mutex создается заново при каждом вызове - не работает!
// 2. mu.Unlock() вызывается дважды для одного Lock() - паника.
// 3. Нет защиты от race condition при параллельных запросах одного ключа.
// 4. Нет оптимизации для параллельного чтения (можно использовать RWMutex).
//
// Дополнительные улучшения для production:
// - Добавить TTL для кэша (expiration).
// - Ограничить размер кэша (LRU eviction).
// - Использовать sync.Map для высококонкурентных сценариев.
// - Добавить метрики (hit rate, miss rate).
