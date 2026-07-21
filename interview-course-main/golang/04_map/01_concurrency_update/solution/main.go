package main

import (
	"fmt"
	"sync"
)

// ConcurrentMap - потокобезопасная мапа.
type ConcurrentMap struct {
	mu   sync.RWMutex
	data map[string]string
}

// NewConcurrentMap создает новую потокобезопасную мапу.
func NewConcurrentMap() *ConcurrentMap {
	return &ConcurrentMap{
		data: make(map[string]string),
	}
}

// GetOrCreate возвращает значение по ключу, создавая его, если не существует.
// Гарантирует, что только одна горутина создаст значение для ключа.
func (cm *ConcurrentMap) GetOrCreate(key, value string) string {
	// Сначала проверяем, существует ли ключ (читаем с RLock для параллельного чтения).
	cm.mu.RLock()
	if val, exists := cm.data[key]; exists {
		cm.mu.RUnlock()
		return val
	}
	cm.mu.RUnlock()

	// Ключа нет, нужно создать. Используем Lock для записи.
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Двойная проверка: возможно, другая горутина уже создала значение.
	if val, exists := cm.data[key]; exists {
		return val
	}

	// Создаем новое значение.
	cm.data[key] = value
	return value
}

func main() {
	cm := NewConcurrentMap()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		val := cm.GetOrCreate("key1", "value1")
		fmt.Println("Goroutine 1 got:", val)
	}()

	go func() {
		defer wg.Done()
		val := cm.GetOrCreate("key1", "value2")
		fmt.Println("Goroutine 2 got:", val)
	}()

	wg.Wait()
}

// Ответ:
// Обе горутины получат одно и то же значение (либо "value1", либо "value2").
// Порядок не гарантирован, но гарантируется, что значение будет создано только один раз.
//
// Объяснение:
// 1. Используем RWMutex для оптимизации: параллельное чтение разрешено.
// 2. Сначала проверяем с RLock (быстрое чтение).
// 3. Если ключа нет, переходим на Lock для записи.
// 4. Двойная проверка (double-checked locking) предотвращает создание дубликатов.
// 5. Только одна горутина создаст значение, остальные получат уже созданное.
