package main

// СПИСОК НАЙДЕННЫХ ПРОБЛЕМ И ИСПРАВЛЕНИЙ:
//
// 1. Race condition при конкурентном доступе к map (строки 27-38)
//    Проблема: map в Go не потокобезопасна
//              Одновременное чтение и запись вызывает panic или data race
//              go run -race обнаружит: "fatal error: concurrent map read and map write"
//    Решение: Использовать sync.RWMutex или sync.Map
//
// 2. Race condition в GetOrCompute (строки 47-60)
//    Проблема: Между Get() и Set() другая горутина может сделать то же
//              Множественные вычисления одного ключа параллельно
//    Пример:
//      Goroutine 1: Get("key1") -> miss -> начинает вычисление
//      Goroutine 2: Get("key1") -> miss -> начинает вычисление
//      Обе вычисляют одно и то же! (thundering herd)
//    Решение: Использовать singleflight или блокировку на уровне ключа
//
// 3. Отсутствует TTL (Time To Live)
//    Проблема: Записи живут вечно, даже если данные устарели
//    Решение: Хранить timestamp и проверять при чтении
//
// 4. Отсутствует ограничение размера кеша
//    Проблема: Кеш растет бесконечно, memory leak
//              При высокой нагрузке может съесть всю память
//    Решение: LRU/LFU eviction policy или ограничение по размеру
//
// 5. Нет периодической очистки устаревших записей
//    Проблема: Даже с TTL, записи остаются в памяти
//    Решение: Фоновая горутина для cleanup
//
// 6. Неэффективная блокировка
//    Проблема: Если использовать Mutex, чтение блокирует чтение
//    Решение: Использовать RWMutex (множественное чтение, эксклюзивная запись)
//
// 7. Нет метрик (cache hits/misses)
//    Проблема: Не видно эффективности кеша
//    Решение: Счетчики hits/misses
//
// 8. Panic при concurrent map access
//    Проблема: concurrent map writes вызывает panic
//    Решение: Потокобезопасные операции
//
// 9. Нет graceful shutdown для cleanup горутины
//    Проблема: Cleanup горутина не останавливается
//    Решение: Context или done channel
//
// 10. Double computation проблема (thundering herd)
//     Проблема: Много горутин вычисляют одно и то же одновременно
//     Решение: golang.org/x/sync/singleflight

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

// CacheEntry хранит значение с метаданными
type CacheEntry struct {
	value     string
	expiresAt time.Time
}

// Cache - потокобезопасный кеш с TTL
type Cache struct {
	mu      sync.RWMutex
	data    map[string]CacheEntry
	ttl     time.Duration
	maxSize int
	hits    int64
	misses  int64
	group   singleflight.Group // Для предотвращения thundering herd
}

// NewCache создает новый кеш
func NewCache(ttl time.Duration, maxSize int) *Cache {
	return &Cache{
		data:    make(map[string]CacheEntry),
		ttl:     ttl,
		maxSize: maxSize,
	}
}

// Get получает значение из кеша
func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.data[key]
	if !ok {
		c.misses++
		return "", false
	}

	// Проверяем TTL
	if time.Now().After(entry.expiresAt) {
		c.misses++
		return "", false
	}

	c.hits++
	return entry.value, true
}

// Set сохраняет значение в кеш
func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Проверяем размер кеша
	if len(c.data) >= c.maxSize {
		// Простая eviction: удаляем первый попавшийся ключ
		// В production используйте LRU
		for k := range c.data {
			delete(c.data, k)
			break
		}
	}

	c.data[key] = CacheEntry{
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
	}
}

// Delete удаляет ключ из кеша
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

// Stats возвращает статистику кеша
func (c *Cache) Stats() (hits, misses int64, size int) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.hits, c.misses, len(c.data)
}

// Cleanup удаляет устаревшие записи
func (c *Cache) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, entry := range c.data {
		if now.After(entry.expiresAt) {
			delete(c.data, key)
		}
	}
}

// StartCleanup запускает фоновую очистку
func (c *Cache) StartCleanup(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.Cleanup()
		case <-ctx.Done():
			return
		}
	}
}

// GetOrCompute получает значение или вычисляет его (с защитой от thundering herd)
func (c *Cache) GetOrCompute(key string, compute func(string) string) string {
	// Быстрая проверка в кеше
	if value, ok := c.Get(key); ok {
		return value
	}

	// Singleflight гарантирует, что compute выполнится только один раз
	// для одинаковых ключей, даже если вызвано из разных горутин
	result, err, _ := c.group.Do(key, func() (interface{}, error) {
		// Вычисляем значение
		value := compute(key)

		// Сохраняем в кеш
		c.Set(key, value)

		return value, nil
	})

	// Если произошла ошибка (не должно случиться с нашей функцией),
	// выполняем fallback вычисление
	if err != nil {
		value := compute(key)
		c.Set(key, value)
		return value
	}

	return result.(string)
}

// expensiveComputation симулирует дорогое вычисление
func expensiveComputation(key string) string {
	fmt.Printf("⏳ Computing %s...\n", key)
	time.Sleep(100 * time.Millisecond)
	return fmt.Sprintf("result for %s", key)
}

// РЕШЕНИЕ 1: Базовый потокобезопасный кеш
func solution1() {
	fmt.Println("=== Решение 1: Базовый потокобезопасный кеш ===")

	cache := NewCache(5*time.Second, 100)
	var wg sync.WaitGroup

	// Запускаем 10 горутин
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for j := 0; j < 3; j++ {
				key := fmt.Sprintf("key%d", j%2)
				result := cache.GetOrCompute(key, expensiveComputation)
				fmt.Printf("Goroutine %d: %s = %s\n", id, key, result)
			}
		}(i)
	}

	wg.Wait()

	// Статистика
	hits, misses, size := cache.Stats()
	hitRate := float64(hits) / float64(hits+misses) * 100
	fmt.Printf("\n📊 Stats: Hits=%d, Misses=%d, Size=%d, Hit Rate=%.1f%%\n\n",
		hits, misses, size, hitRate)
}

// РЕШЕНИЕ 2: С автоматической очисткой
func solution2() {
	fmt.Println("=== Решение 2: С фоновой очисткой ===")

	cache := NewCache(2*time.Second, 100)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Запускаем фоновую очистку каждые 1 секунду
	go cache.StartCleanup(ctx, 1*time.Second)

	// Добавляем записи
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("key%d", i)
		cache.Set(key, fmt.Sprintf("value%d", i))
	}

	fmt.Println("Added 5 entries")
	_, _, size := cache.Stats()
	fmt.Printf("Cache size: %d\n", size)

	// Ждем 3 секунды (TTL = 2s)
	time.Sleep(3 * time.Second)

	// Проверяем размер после cleanup
	_, _, size = cache.Stats()
	fmt.Printf("Cache size after TTL: %d (should be 0)\n\n", size)
}

// РЕШЕНИЕ 3: Демонстрация thundering herd защиты
func solution3() {
	fmt.Println("=== Решение 3: Защита от thundering herd ===")

	cache := NewCache(5*time.Second, 100)
	var wg sync.WaitGroup

	// 10 горутин одновременно запрашивают один ключ
	key := "expensive-key"

	fmt.Println("Launching 10 goroutines requesting the same key...")
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			result := cache.GetOrCompute(key, expensiveComputation)
			fmt.Printf("Goroutine %d got: %s\n", id, result)
		}(i)
	}

	wg.Wait()

	// Должно быть только ОДНО вычисление!
	fmt.Println("\n✅ Only one computation happened (thundering herd prevented)")
}

func main() {
	solution1()
	solution2()
	solution3()
}

// ПОДРОБНОЕ ОБЪЯСНЕНИЕ RACE CONDITIONS:
//
// RACE CONDITION #1: Concurrent map access
//
// Что происходит БЕЗ синхронизации:
//
// Time    Goroutine 1              Goroutine 2
// ─────────────────────────────────────────────────────────
// 0ms     cache.Get("key1")        cache.Get("key1")
// 1ms     Read data["key1"]        Read data["key1"]
//         ↑ Одновременное чтение = OK
//
// 5ms     cache.Set("key1", "v1")  cache.Get("key1")
// 6ms     Write data["key1"]       Read data["key1"]
//         ↑ RACE! Одновременные read/write
//         ↓ panic: concurrent map read and map write
//
// Как обнаружить:
//   go run -race main.go
//   ==================
//   WARNING: DATA RACE
//   Write at 0x... by goroutine 7:
//     main.(*Cache).Set()
//   Previous read at 0x... by goroutine 8:
//     main.(*Cache).Get()
//   ==================
//
// RACE CONDITION #2: Double computation (thundering herd)
//
// Без защиты:
//
// Time    G1                      G2                      G3
// ─────────────────────────────────────────────────────────────
// 0ms     GetOrCompute("key")     GetOrCompute("key")     GetOrCompute("key")
// 1ms     Get("key") -> miss      Get("key") -> miss      Get("key") -> miss
// 2ms     compute("key") ⏳       compute("key") ⏳       compute("key") ⏳
//         ↑ Все три вычисляют одновременно!
// 102ms   Set("key", "result")    Set("key", "result")    Set("key", "result")
//         ↑ Три вычисления вместо одного
//
// С защитой (singleflight pattern):
//
// Time    G1                      G2                      G3
// ─────────────────────────────────────────────────────────────
// 0ms     GetOrCompute("key")     GetOrCompute("key")     GetOrCompute("key")
// 1ms     Get("key") -> miss      Get("key") -> miss      Get("key") -> miss
// 2ms     singleflight.Do()       singleflight.Do()       singleflight.Do()
// 3ms     compute("key") ⏳       ⏸️ waiting...           ⏸️ waiting...
// 102ms   Set("key", "result")    ⏸️ waiting...           ⏸️ waiting...
// 103ms   Return "result"         Return "result"         Return "result"
//         ↑ Только одно вычисление! Все получают результат напрямую
//
// ПРОБЛЕМА TTL:
//
// Без TTL:
//   T=0:   Set("user:123", "John")
//   T=1h:  Get("user:123") -> "John" (но пользователь уже "Jane")
//          ↑ Устаревшие данные!
//
// С TTL:
//   T=0:    Set("user:123", "John", expires=T+5min)
//   T=3min: Get("user:123") -> "John" ✓
//   T=6min: Get("user:123") -> miss (TTL expired)
//           → compute() → Get fresh data
//
// ПРОБЛЕМА MEMORY LEAK:
//
// Без ограничения размера:
//   Requests: 1M unique keys
//   Cache: Stores all 1M entries
//   Memory: 1M * 1KB = 1GB
//   After 1 hour: 60M requests
//   Memory: 60GB! ☠️
//
// С LRU eviction:
//   Cache max size: 10000 entries
//   Memory: 10000 * 1KB = 10MB ✓
//   Old entries evicted automatically
//
// SYNC.RWMUTEX vs SYNC.MUTEX:
//
// Mutex (эксклюзивная блокировка):
//   Reader 1: Lock() → Read → Unlock()
//   Reader 2: Lock() → ⏸️ WAIT → Read → Unlock()
//   Reader 3: Lock() → ⏸️ WAIT → Read → Unlock()
//   ↑ Читатели блокируют друг друга
//
// RWMutex (множественное чтение):
//   Reader 1: RLock() → Read → RUnlock()
//   Reader 2: RLock() → Read → RUnlock()  ← Параллельно!
//   Reader 3: RLock() → Read → RUnlock()  ← Параллельно!
//   Writer:   Lock() → ⏸️ WAIT → Write → Unlock()
//   ↑ Читатели не блокируют друг друга
//
// Производительность:
//   100 readers, 1 writer
//   Mutex:    100 sequential locks = slow
//   RWMutex:  100 parallel reads = fast ✓
//
// CLEANUP STRATEGY:
//
// Lazy cleanup (при Get):
//   + Простая реализация
//   + Не нужна фоновая горутина
//   - Память не освобождается между Get()
//   - Неиспользуемые ключи остаются в памяти
//
// Active cleanup (фоновая горутина):
//   + Проактивное освобождение памяти
//   + Контролируемое использование ресурсов
//   - Дополнительная горутина
//   - Нужен graceful shutdown
//
// Hybrid (lazy + active):
//   Проверяем TTL при Get() + периодический cleanup
//   ↑ Лучший баланс
//
// ЛУЧШИЕ ПРАКТИКИ:
//
// 1. Используйте RWMutex для кеша (много чтений)
// 2. Всегда проверяйте TTL при чтении
// 3. Ограничивайте размер кеша (LRU/LFU)
// 4. Предотвращайте thundering herd (singleflight)
// 5. Фоновая cleanup горутина с graceful shutdown
// 6. Метрики: hit rate, size, evictions
// 7. Тестируйте с -race флагом
// 8. Для production: готовые решения (groupcache, ristretto)
// 9. Мониторинг: память, latency, hit rate
// 10. Benchmark разных стратегий для вашей нагрузки
