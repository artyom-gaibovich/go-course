package main

import (
	"fmt"
	"sync"
)

type Cache struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

func (c *Cache) Store(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

func (c *Cache) Load(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.data[key]
	return value, ok
}

func main() {
	cache := &Cache{
		data: make(map[string]interface{}),
	}

	cache.Store("name", "Alice")
	cache.Store("age", 25)

	// ИСПРАВЛЕНИЕ: безопасное приведение типов с проверкой
	name, ok := cache.Load("name")
	if ok {
		if nameStr, ok := name.(string); ok {
			fmt.Println("Name:", nameStr)
		}
	}

	age, ok := cache.Load("age")
	if ok {
		if ageInt, ok := age.(int); ok {
			fmt.Println("Age:", ageInt)
		}
	}

	height, ok := cache.Load("height")
	if !ok || height == nil {
		fmt.Println("Height not found")
	}

	width, ok := cache.Load("width")
	if ok {
		if widthFloat, ok := width.(float64); ok {
			fmt.Println("Width:", widthFloat)
		} else {
			fmt.Println("Width not found or wrong type")
		}
	} else {
		fmt.Println("Width not found")
	}
}

// Проблемы в исходном коде:
// 1. Type assertion без проверки: cache.Load("width").(float64) вызовет panic,
//    если ключа нет или тип не соответствует.
// 2. Load не возвращает информацию о наличии ключа (как map[key]).
// 3. Нет проверки результата type assertion.
//
// Решения:
// 1. Изменить сигнатуру Load для возврата (value, ok).
// 2. Использовать безопасное приведение типов с проверкой.
// 3. Проверять наличие ключа перед type assertion.
