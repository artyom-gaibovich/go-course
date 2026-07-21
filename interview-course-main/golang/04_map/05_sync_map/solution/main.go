package main

import (
	"fmt"
	"sync"
	"time"
)

// Cache - потокобезопасный кэш с вычислением значений на основе sync.Map.
type Cache struct {
	data sync.Map
}

// NewCache создает новый кэш.
func NewCache() *Cache {
	return &Cache{}
}

// GetOrCompute возвращает значение из кэша или вычисляет его.
// Гарантирует, что функция compute выполнится только один раз для каждого ключа.
func (c *Cache) GetOrCompute(key string, compute func() interface{}) interface{} {
	// Пытаемся загрузить значение из кэша.
	if val, ok := c.data.Load(key); ok {
		return val
	}

	// Значения нет в кэше, нужно вычислить.
	// Создаем канал для синхронизации вычисления между горутинами.
	ch := make(chan interface{}, 1)

	// Пытаемся сохранить канал. Если другая горутина уже создала канал, получим его.
	actual, loaded := c.data.LoadOrStore(key, ch)

	if loaded {
		// Другая горутина уже вычисляет значение, ждем результат из канала.
		return <-actual.(chan interface{})
	}

	// Мы первые, кто начал вычисление для этого ключа.
	// Вычисляем значение.
	result := compute()

	// Сохраняем результат вместо канала.
	c.data.Store(key, result)

	// Отправляем результат в канал для других горутин, которые ждут.
	ch <- result
	close(ch)

	return result
}

func main() {
	cache := NewCache()

	// Симулируем дорогостоящую операцию.
	expensiveOperation := func() interface{} {
		fmt.Println("Выполняю дорогостоящую операцию...")
		time.Sleep(1 * time.Second)
		return "bla-bla"
	}

	// Первый вызов - вычисление.
	result1 := cache.GetOrCompute("key1", expensiveOperation)
	// Конвертируем из interface{} в string
	if str, ok := result1.(string); ok {
		fmt.Println("Результат 1:", str)
	}

	// Второй вызов - из кэша (дорогостоящая операция не выполнится).
	result2 := cache.GetOrCompute("key1", expensiveOperation)
	// Конвертируем из interface{} в string
	if str, ok := result2.(string); ok {
		fmt.Println("Результат 2:", str)
	}

	// Параллельные вызовы - только один выполнит вычисление.
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			result := cache.GetOrCompute("key2", func() interface{} {
				fmt.Printf("Горутина %d выполняет вычисление\n", id)
				time.Sleep(500 * time.Millisecond)
				return fmt.Sprintf("результат от горутины %d", id)
			})
			// Конвертируем из interface{} в string
			if str, ok := result.(string); ok {
				fmt.Printf("Горутина %d получила: %s\n", id, str)
			}
		}(i)
	}
	wg.Wait()
}

// Объяснение:
// 1. Используем sync.Map для потокобезопасного хранения данных.
// 2. LoadOrStore с каналом гарантирует, что compute выполнится только один раз.
// 3. sync.Map оптимизирован для случаев с большим количеством чтений и редкими записями.
// 4. Все горутины получат одно и то же вычисленное значение.
//
// Как это работает:
// - Первая горутина создает канал и сохраняет его через LoadOrStore.
// - Остальные горутины получают этот канал и ждут результата.
// - Первая горутина вычисляет значение, сохраняет его и отправляет в канал.
// - Все горутины получают результат и возвращают его.
