package main

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"
)

var (
	counter int64
	mu      sync.Mutex
)

// SimulateRequest выполняет запрос с таймаутом и возвращает значение счетчика.
func SimulateRequest(ctx context.Context) (int64, error) {
	start := time.Now()

	// Симулируем случайную задержку
	delay := time.Duration(rand.Int63n(5)) * time.Second

	select {
	case <-time.After(delay):
		// Запрос выполнен успешно
		mu.Lock()
		counter++
		val := counter
		mu.Unlock()

		elapsed := time.Since(start)
		log.Printf("Запрос выполнен за %v, значение счетчика: %d\n", elapsed, val)
		return val, nil

	case <-ctx.Done():
		// Запрос отменен по таймауту
		elapsed := time.Since(start)
		log.Printf("Запрос отменен по таймауту через %v\n", elapsed)
		return 0, ctx.Err()
	}
}

func main() {
	// Создаем контекст с таймаутом 3 секунды
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	val, err := SimulateRequest(ctx)
	if err != nil {
		log.Printf("Ошибка: %v\n", err)
	} else {
		log.Printf("Значение счетчика: %d\n", val)
	}
}

// Объяснение:
// 1. Используем context.WithTimeout для установки таймаута.
// 2. sync.Mutex для защиты counter от race condition.
// 3. select для ожидания либо завершения запроса, либо таймаута.
// 4. Логируем время выполнения запроса.
