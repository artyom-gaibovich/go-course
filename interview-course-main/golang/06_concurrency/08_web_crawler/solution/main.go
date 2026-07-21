package main

import (
	"fmt"
	"sync"
	"time"
)

// Solution интерфейс для статистики web crawler.
type Solution interface {
	OnPageDownloaded(host string, timestamp time.Time)
	Count(host string) int
}

// ============================================================================
// Решение 1: Хранение временных меток (простое, но требует больше памяти)
// ============================================================================

type crawlerStats struct {
	mu        sync.RWMutex           // Защита от конкурентного доступа
	downloads map[string][]time.Time // host -> список временных меток
	window    time.Duration          // Временное окно (10 минут)
}

// NewCrawlerStats создает новый экземпляр статистики crawler.
func NewCrawlerStats() Solution {
	return &crawlerStats{
		downloads: make(map[string][]time.Time),
		window:    10 * time.Minute,
	}
}

// OnPageDownloaded регистрирует загрузку страницы.
// Метод потокобезопасен.
func (c *crawlerStats) OnPageDownloaded(host string, timestamp time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Добавляем временную метку загрузки.
	c.downloads[host] = append(c.downloads[host], timestamp)
}

// Count возвращает количество загрузок хоста за последние 10 минут.
// Метод потокобезопасен.
func (c *crawlerStats) Count(host string) int {
	c.mu.Lock()
	defer c.mu.Unlock()

	timestamps, exists := c.downloads[host]
	if !exists {
		return 0
	}

	// Текущее время для вычисления окна.
	now := time.Now()
	windowStart := now.Add(-c.window)

	// Находим индекс первого валидного элемента.
	validStart := 0
	for i, ts := range timestamps {
		if !ts.Before(windowStart) {
			validStart = i
			break
		}
	}

	// Обрезаем устаревшие данные для экономии памяти.
	validTimestamps := timestamps[validStart:]
	c.downloads[host] = validTimestamps

	return len(validTimestamps)
}

// ============================================================================
// Решение 2: Bucket approach (меньше памяти, чуть менее точный)
// ============================================================================

// Группируем загрузки по временным bucket'ам (например, по минутам).
// Вместо хранения каждой временной метки храним только счетчики.

type bucketCrawlerStats struct {
	mu         sync.RWMutex
	buckets    map[string]map[int64]int // host -> bucket_timestamp -> count
	window     time.Duration
	bucketSize time.Duration // Размер bucket (например, 1 минута)
}

// NewBucketCrawlerStats создает статистику с bucket подходом.
func NewBucketCrawlerStats(bucketSize time.Duration) Solution {
	return &bucketCrawlerStats{
		buckets:    make(map[string]map[int64]int),
		window:     10 * time.Minute,
		bucketSize: bucketSize,
	}
}

// OnPageDownloaded увеличивает счетчик для соответствующего bucket.
func (c *bucketCrawlerStats) OnPageDownloaded(host string, timestamp time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Вычисляем bucket: округляем timestamp до bucketSize.
	// ЗДЕСЬ ИСПОЛЬЗУЕТСЯ time.Truncate() - см. подробное объяснение ниже.
	bucket := timestamp.Truncate(c.bucketSize).Unix()

	// Инициализируем map для хоста, если нужно.
	if c.buckets[host] == nil {
		c.buckets[host] = make(map[int64]int)
	}

	// Увеличиваем счетчик для bucket.
	c.buckets[host][bucket]++
}

// Count суммирует счетчики из bucket'ов за последние 10 минут.
func (c *bucketCrawlerStats) Count(host string) int {
	c.mu.Lock()
	defer c.mu.Unlock()

	hostBuckets, exists := c.buckets[host]
	if !exists {
		return 0
	}

	// Вычисляем границу временного окна.
	now := time.Now()
	windowStart := now.Add(-c.window)
	// СНОВА используем Truncate() для округления
	windowStartBucket := windowStart.Truncate(c.bucketSize).Unix()

	// Суммируем счетчики из валидных bucket'ов.
	count := 0
	for bucket, bucketCount := range hostBuckets {
		if bucket >= windowStartBucket {
			count += bucketCount
		} else {
			// Удаляем устаревшие bucket'ы для экономии памяти.
			delete(hostBuckets, bucket)
		}
	}

	return count
}

func main() {
	// Демонстрация двух подходов

	fmt.Println("=== Решение 1: Хранение временных меток ===")
	stats1 := NewCrawlerStats()
	demonstrateStats(stats1)

	fmt.Println("\n=== Решение 2: Bucket approach ===")
	stats2 := NewBucketCrawlerStats(1 * time.Minute)
	demonstrateStats(stats2)

	fmt.Println("\n=== Демонстрация time.Truncate() ===")
	demonstrateTruncate()
}

func demonstrateStats(stats Solution) {
	now := time.Now()

	// Симулируем загрузки страниц
	stats.OnPageDownloaded("facebook.com", now.Add(-11*time.Minute)) // Старая
	stats.OnPageDownloaded("facebook.com", now.Add(-9*time.Minute))
	stats.OnPageDownloaded("facebook.com", now.Add(-5*time.Minute))
	stats.OnPageDownloaded("facebook.com", now.Add(-2*time.Minute))
	stats.OnPageDownloaded("facebook.com", now.Add(-1*time.Minute))

	stats.OnPageDownloaded("twitter.com", now.Add(-8*time.Minute))
	stats.OnPageDownloaded("twitter.com", now.Add(-3*time.Minute))

	// Получаем статистику
	fmt.Printf("facebook.com: %d\n", stats.Count("facebook.com")) // 4
	fmt.Printf("twitter.com: %d\n", stats.Count("twitter.com"))   // 2
	fmt.Printf("unknown.com: %d\n", stats.Count("unknown.com"))   // 0
}

func demonstrateTruncate() {
	// Примеры работы time.Truncate()

	// Пример 1: Округление до минуты
	t1 := time.Date(2024, 1, 15, 14, 23, 47, 123456789, time.UTC)
	fmt.Printf("Исходное время: %s\n", t1.Format("15:04:05.000000000"))
	fmt.Printf("Truncate(1m):    %s\n", t1.Truncate(time.Minute).Format("15:04:05.000000000"))
	fmt.Printf("Truncate(10m):   %s\n", t1.Truncate(10*time.Minute).Format("15:04:05.000000000"))
	fmt.Printf("Truncate(1h):    %s\n", t1.Truncate(time.Hour).Format("15:04:05.000000000"))

	// Пример 2: Bucket для статистики
	fmt.Println("\nПример bucket'ов:")
	events := []time.Time{
		time.Date(2024, 1, 15, 14, 23, 10, 0, time.UTC),
		time.Date(2024, 1, 15, 14, 23, 45, 0, time.UTC),
		time.Date(2024, 1, 15, 14, 24, 30, 0, time.UTC),
		time.Date(2024, 1, 15, 14, 25, 15, 0, time.UTC),
	}

	fmt.Println("События группируются в один bucket:")
	for _, event := range events {
		bucket := event.Truncate(time.Minute)
		fmt.Printf("  Событие: %s → Bucket: %s\n",
			event.Format("15:04:05"),
			bucket.Format("15:04:00"))
	}
}

// Объяснение решений:
//
// ============================================================================
// ЧТО ТАКОЕ time.Truncate() И КАК ОНО РАБОТАЕТ
// ============================================================================
//
// time.Truncate(d Duration) округляет время ВНИЗ до ближайшего кратного d.
// Это как математический floor(), но для временных меток.
//
// СИНТАКСИС:
//   truncated := timestamp.Truncate(duration)
//
// ПРИМЕРЫ:
//
// Пример 1: Округление до минуты
//   t := 2024-01-15 14:23:47.123456789
//   t.Truncate(time.Minute) = 2024-01-15 14:23:00.000000000
//   Отбросили: 47 секунд и наносекунды
//
// Пример 2: Округление до 10 минут
//   t := 2024-01-15 14:23:47
//   t.Truncate(10 * time.Minute) = 2024-01-15 14:20:00
//   Округлили вниз к ближайшему 10-минутному интервалу
//
// Пример 3: Округление до часа
//   t := 2024-01-15 14:23:47
//   t.Truncate(time.Hour) = 2024-01-15 14:00:00
//   Отбросили минуты, секунды, наносекунды
//
// Пример 4: Округление до дня (24 часа)
//   t := 2024-01-15 14:23:47
//   t.Truncate(24 * time.Hour) = 2024-01-15 00:00:00
//   Получили начало дня
//
// КАК РАБОТАЕТ ВНУТРИ:
//
// Алгоритм Truncate() можно представить так:
//
//   func (t Time) Truncate(d Duration) Time {
//       if d <= 0 {
//           return t
//       }
//       // Получаем Unix timestamp в наносекундах
//       ns := t.UnixNano()
//       // Вычисляем remainder (остаток от деления)
//       r := ns % int64(d)
//       // Вычитаем remainder, чтобы округлить вниз
//       return time.Unix(0, ns - r)
//   }
//
// Пошагово для t = 14:23:47.5, d = 1 minute:
//   1. ns = Unix time в наносекундах
//   2. d = 1 minute = 60_000_000_000 наносекунд
//   3. r = ns % d = остаток (47.5 секунды в наносекундах)
//   4. ns - r = убираем остаток, получаем 14:23:00.0
//
// ПРИМЕНЕНИЕ В BUCKET APPROACH:
//
// Идея: группировать события по временным интервалам (bucket'ам).
//
// Без Truncate (сохраняем каждое событие):
//   14:23:10 → сохранить timestamp
//   14:23:45 → сохранить timestamp
//   14:24:30 → сохранить timestamp
//   Результат: 3 записи в памяти
//
// С Truncate (группируем в bucket'ы):
//   14:23:10 → Truncate(1m) → 14:23:00 → count[14:23:00]++
//   14:23:45 → Truncate(1m) → 14:23:00 → count[14:23:00]++
//   14:24:30 → Truncate(1m) → 14:24:00 → count[14:24:00]++
//   Результат: 2 bucket'а (14:23:00 → 2, 14:24:00 → 1)
//
// Преимущества:
// - Меньше памяти: храним счетчики, а не все timestamps
// - Быстрее Count(): обходим bucket'ы вместо всех событий
// - Проще очистка: удаляем целые bucket'ы
//
// ВЫБОР РАЗМЕРА BUCKET:
//
// Размер bucket влияет на точность и память:
//
// bucket = 1 секунда:
//   Плюсы: высокая точность (±1 сек)
//   Минусы: больше bucket'ов, больше памяти
//   Для окна 10 минут: 600 bucket'ов максимум
//
// bucket = 1 минута:
//   Плюсы: баланс точности и памяти
//   Минусы: точность ±1 минута
//   Для окна 10 минут: 10 bucket'ов максимум
//
// bucket = 10 минут:
//   Плюсы: минимум памяти
//   Минусы: низкая точность (±10 минут)
//   Для окна 10 минут: 1 bucket
//
// Правило выбора:
//   bucketSize = window / desired_granularity
//   Например: window=10m, granularity=60 → bucketSize=10s
//
// СРАВНЕНИЕ С ДРУГИМИ ОПЕРАЦИЯМИ:
//
// time.Round() vs time.Truncate():
//
//   t := 14:23:47
//
//   Truncate(1m):  14:23:00  (округление вниз)
//   Round(1m):     14:24:00  (округление к ближайшему)
//
//   Truncate(10m): 14:20:00  (округление вниз)
//   Round(10m):    14:20:00  (ближайшее — совпало)
//
// Когда использовать:
// - Truncate: для bucket'ов, статистики, group by
// - Round: для отображения времени пользователю
//
// time.Add() vs time.Truncate():
//
//   t := 14:23:47
//
//   Add(-47*Second):   14:23:00  (точное вычитание)
//   Truncate(1m):      14:23:00  (округление вниз)
//
// Разница:
// - Add требует знать точное смещение
// - Truncate работает для любого времени в минуте
//
// ПРОИЗВОДИТЕЛЬНОСТЬ:
//
// Truncate() очень быстрая операция:
// - Простая арифметика: ns % d и ns - remainder
// - Нет аллокаций
// - O(1) сложность
//
// ПОДВОДНЫЕ КАМНИ:
//
// 1. Часовые пояса:
//    Truncate работает с monotonic time, игнорирует timezone.
//    Для группировки по дням используйте:
//      y, m, d := t.Date()
//      dayStart := time.Date(y, m, d, 0, 0, 0, 0, t.Location())
//
// 2. Не кратные интервалы:
//    Truncate(25*time.Minute) может дать неожиданные результаты.
//    25 минут не делит час нацело.
//
// 3. Нулевая и отрицательная duration:
//    t.Truncate(0) вернет t без изменений
//    t.Truncate(-1) вернет t без изменений
//
// ПРИМЕРЫ ИСПОЛЬЗОВАНИЯ В PRODUCTION:
//
// 1. Метрики (Prometheus, StatsD):
//    Группируют метрики в bucket'ы для агрегации
//    bucket := timestamp.Truncate(10 * time.Second)
//
// 2. Логирование:
//    Ротация логов по часам
//    hour := time.Now().Truncate(time.Hour)
//    logFile := fmt.Sprintf("app-%s.log", hour.Format("2006-01-02-15"))
//
// 3. Rate limiting:
//    Окна для sliding window rate limiter
//    window := time.Now().Truncate(time.Minute)
//
// 4. Caching:
//    Cache key с минутной точностью
//    cacheKey := fmt.Sprintf("data:%d", time.Now().Truncate(time.Minute).Unix())
//
// 5. Базы данных (Time-series DB):
//    Partitioning по временным bucket'ам
//    partition := timestamp.Truncate(24 * time.Hour) // по дням
//
// ============================================================================
// Решение 1: Хранение временных меток
// ============================================================================
//
// Плюсы:
// - Простая реализация
// - Точное временное окно
// - Легко понять и отладить
//
// Минусы:
// - Память растет линейно с количеством загрузок
// - O(n) для очистки старых данных
// - Не подходит для очень высоких нагрузок
//
// Когда использовать:
// - Умеренное количество загрузок (< 100K/мин на хост)
// - Нужна точность
// - Простота важнее оптимизации
//
// ============================================================================
// Решение 2: Bucket approach
// ============================================================================
//
// Плюсы:
// - Меньше памяти (храним счетчики, а не все timestamps)
// - Быстрая очистка (удаляем целые bucket'ы)
// - Масштабируется лучше для высоких нагрузок
// - ИСПОЛЬЗУЕТ time.Truncate() для группировки
//
// Минусы:
// - Немного менее точное временное окно (зависит от размера bucket)
// - Чуть сложнее реализация
//
// Когда использовать:
// - Высокие нагрузки (> 100K/мин на хост)
// - Можно пожертвовать точностью до размера bucket
// - Важна память
//
// Пример с разными bucket sizes:
//
// bucket = 1 секунда, window = 10 минут:
//   Максимум bucket'ов: 600
//   Точность: ±1 секунда
//   Память на хост: ~600 × (8 bytes key + 8 bytes value) = 9.6 KB
//
// bucket = 1 минута, window = 10 минут:
//   Максимум bucket'ов: 10
//   Точность: ±1 минута
//   Память на хост: ~10 × 16 bytes = 160 bytes
//
// ============================================================================
// Сравнение производительности
// ============================================================================
//
// Память (для 1M загрузок/мин за 10 мин):
// - Решение 1: ~10M timestamps × 24 bytes = 240 MB на хост
// - Решение 2: ~600 buckets × 16 bytes = 10 KB на хост (bucket=1 сек)
//
// Разница: 24000x !
//
// Сложность Count():
// - Решение 1: O(n) где n = количество записей (до 10M)
// - Решение 2: O(b) где b = количество bucket'ов (~600)
//
// Сложность OnPageDownloaded():
// - Решение 1: O(1) append
// - Решение 2: O(1) increment + Truncate O(1)
//
// ============================================================================
// Рекомендации для собеседования
// ============================================================================
//
// 1. Начните с простого решения (Решение 1)
// 2. Обсудите trade-offs и ограничения
// 3. Спросите о требованиях:
//    - Ожидаемый RPS?
//    - Количество хостов?
//    - Ограничения памяти?
//    - Допустимая погрешность?
// 4. Предложите оптимизацию через bucket approach
// 5. Объясните роль Truncate() в группировке
// 6. Обсудите масштабирование:
//    - Sharding по хостам
//    - Распределенная статистика (Redis, ClickHouse)
//    - Lock-free структуры для высокой конкурентности
//
// ============================================================================
// Дополнительные оптимизации
// ============================================================================
//
// 1. Фоновая очистка:
//    ticker := time.NewTicker(1 * time.Minute)
//    go func() {
//        for range ticker.C {
//            cleanupOldBuckets()
//        }
//    }()
//
// 2. Sharding:
//    // Разделяем хосты по 16 shard'ам
//    shardID := hash(host) % 16
//    stats[shardID].OnPageDownloaded(host, ts)
//
// 3. Atomic counters:
//    // Для bucket'ов можно использовать atomic.AddInt64
//    atomic.AddInt64(&buckets[key], 1)
//
// 4. Lazy cleanup:
//    // Очищаем только при Count(), а не при каждом OnPageDownloaded
//
// ============================================================================
// КЛЮЧЕВЫЕ ВЫВОДЫ
// ============================================================================
//
// - time.Truncate() округляет время ВНИЗ до кратного duration
// - Идеально для группировки событий в bucket'ы
// - Bucket approach снижает память в тысячи раз
// - Trade-off: точность vs память
// - Выбор bucketSize зависит от требований
// - O(1) операция, очень быстрая
// - Используется в метриках, логах, rate limiting, caching
