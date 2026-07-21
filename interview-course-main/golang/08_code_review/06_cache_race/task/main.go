// Задача: Сделайте ревью кода и исправьте проблемы
//
// ТЗ: простой in-memory кеш для хранения результатов дорогих вычислений.
// Программа кеширует результаты и переиспользует их при повторных запросах.
//
// Этот код НАМЕРЕННО содержит ошибки для учебных целей!
// Не запускайте в production!

package main

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	data map[string]string
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]string),
	}
}

func (c *Cache) Get(key string) (string, bool) {
	value, ok := c.data[key]
	return value, ok
}

func (c *Cache) Set(key, value string) {
	c.data[key] = value
}

func (c *Cache) Delete(key string) {
	delete(c.data, key)
}

// expensiveComputation симулирует дорогое вычисление
func expensiveComputation(key string) string {
	time.Sleep(100 * time.Millisecond)
	return fmt.Sprintf("result for %s", key)
}

// GetOrCompute получает значение из кеша или вычисляет его
func GetOrCompute(cache *Cache, key string) string {
	// Проверяем кеш
	if value, ok := cache.Get(key); ok {
		return value
	}

	// Вычисляем значение
	value := expensiveComputation(key)

	// Сохраняем в кеш
	cache.Set(key, value)

	return value
}

func main() {
	cache := NewCache()
	var wg sync.WaitGroup

	// Запускаем 10 горутин, которые обращаются к кешу
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Каждая горутина делает несколько запросов
			for j := 0; j < 5; j++ {
				key := fmt.Sprintf("key%d", j%3)
				result := GetOrCompute(cache, key)
				fmt.Printf("Goroutine %d: %s = %s\n", id, key, result)
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("Done")
}
