package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

// ============================================================================
// Интерфейсы и базовые типы
// ============================================================================

type Request interface{}

type Response interface{}

type Backend interface {
	Invoke(ctx context.Context, req Request) (Response, error)
}

// BackendImpl реализует Backend для одного экземпляра микросервиса.
type BackendImpl struct {
	addr string
}

var _ Backend = &BackendImpl{}

// NewBackend создает backend для конкретного экземпляра по адресу.
func NewBackend(addr string) *BackendImpl {
	return &BackendImpl{addr: addr}
}

// Invoke отправляет запрос на конкретный backend.
// В реальности здесь был бы HTTP/gRPC вызов.
func (b *BackendImpl) Invoke(ctx context.Context, req Request) (Response, error) {
	// Симуляция обработки запроса
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return fmt.Sprintf("response from %s for request %v", b.addr, req), nil
	}
}

// ============================================================================
// Решение 1: Простой Round Robin с sync.Mutex
// ============================================================================

type Balancer struct {
	backends []Backend // Список backend'ов
	mu       sync.Mutex
	next     int // Индекс следующего backend'а
}

var _ Backend = &Balancer{}

// NewBalancer создает балансировщик для списка адресов backend'ов.
func NewBalancer(addrs []string) *Balancer {
	// Создаем Backend для каждого адреса.
	backends := make([]Backend, len(addrs))
	for i, addr := range addrs {
		backends[i] = NewBackend(addr)
	}

	return &Balancer{
		backends: backends,
		next:     0,
	}
}

// Invoke распределяет запросы между backend'ами по Round Robin.
// Алгоритм Round Robin: берем backend'ы по кругу, один за другим.
func (bal *Balancer) Invoke(ctx context.Context, req Request) (Response, error) {
	if len(bal.backends) == 0 {
		return nil, errors.New("no backends available")
	}

	// Получаем следующий backend потокобезопасно.
	bal.mu.Lock()
	backend := bal.backends[bal.next]
	// Переходим к следующему backend'у по кругу.
	bal.next = (bal.next + 1) % len(bal.backends)
	bal.mu.Unlock()

	// Отправляем запрос на выбранный backend.
	return backend.Invoke(ctx, req)
}

// ============================================================================
// Решение 2: Round Robin с atomic для лучшей производительности
// ============================================================================

type AtomicBalancer struct {
	backends []Backend
	counter  uint64 // Атомарный счетчик для Round Robin
}

var _ Backend = &AtomicBalancer{}

// NewAtomicBalancer создает балансировщик с atomic операциями.
func NewAtomicBalancer(addrs []string) *AtomicBalancer {
	backends := make([]Backend, len(addrs))
	for i, addr := range addrs {
		backends[i] = NewBackend(addr)
	}

	return &AtomicBalancer{
		backends: backends,
		counter:  0,
	}
}

// Invoke использует atomic.AddUint64 для потокобезопасного счетчика.
func (bal *AtomicBalancer) Invoke(ctx context.Context, req Request) (Response, error) {
	if len(bal.backends) == 0 {
		return nil, errors.New("no backends available")
	}

	// Атомарно увеличиваем счетчик и получаем индекс.
	// atomic.AddUint64 гарантирует уникальное значение для каждой горутины.
	idx := atomic.AddUint64(&bal.counter, 1)

	// Вычисляем индекс backend'а по модулю.
	backend := bal.backends[idx%uint64(len(bal.backends))]

	return backend.Invoke(ctx, req)
}

// ============================================================================
// Решение 3: Round Robin с retry и health checks
// ============================================================================

type SmartBalancer struct {
	backends   []Backend
	counter    uint64
	maxRetries int
}

var _ Backend = &SmartBalancer{}

// NewSmartBalancer создает балансировщик с поддержкой retry.
func NewSmartBalancer(addrs []string, maxRetries int) *SmartBalancer {
	backends := make([]Backend, len(addrs))
	for i, addr := range addrs {
		backends[i] = NewBackend(addr)
	}

	return &SmartBalancer{
		backends:   backends,
		counter:    0,
		maxRetries: maxRetries,
	}
}

// Invoke пытается выполнить запрос, при ошибке переходит к следующему backend'у.
func (bal *SmartBalancer) Invoke(ctx context.Context, req Request) (Response, error) {
	if len(bal.backends) == 0 {
		return nil, errors.New("no backends available")
	}

	// Начальный индекс для Round Robin.
	startIdx := atomic.AddUint64(&bal.counter, 1)

	var lastErr error

	// Пробуем до maxRetries backend'ов.
	for i := 0; i < bal.maxRetries; i++ {
		// Вычисляем индекс текущего backend'а.
		idx := (startIdx + uint64(i)) % uint64(len(bal.backends))
		backend := bal.backends[idx]

		// Проверяем context перед вызовом.
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		// Пытаемся выполнить запрос.
		resp, err := backend.Invoke(ctx, req)
		if err == nil {
			return resp, nil
		}

		lastErr = err
		// Можно добавить логирование или метрики здесь
	}

	return nil, fmt.Errorf("all backends failed, last error: %w", lastErr)
}

// ============================================================================
// Демонстрация и тестирование
// ============================================================================

func main() {
	fmt.Println("=== Решение 1: Простой Round Robin ===")
	testBalancer1()

	fmt.Println("\n=== Решение 2: Atomic Round Robin ===")
	testBalancer2()

	fmt.Println("\n=== Решение 3: Smart Balancer с retry ===")
	testBalancer3()

	fmt.Println("\n=== Тест конкурентности ===")
	testConcurrency()
}

func testBalancer1() {
	addrs := []string{"backend1:8080", "backend2:8080", "backend3:8080"}
	balancer := NewBalancer(addrs)

	ctx := context.Background()

	// Отправляем 6 запросов и смотрим распределение.
	for i := 0; i < 6; i++ {
		resp, err := balancer.Invoke(ctx, fmt.Sprintf("request-%d", i))
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		fmt.Printf("Request %d: %v\n", i, resp)
	}
}

func testBalancer2() {
	addrs := []string{"backend1:8080", "backend2:8080", "backend3:8080"}
	balancer := NewAtomicBalancer(addrs)

	ctx := context.Background()

	for i := 0; i < 6; i++ {
		resp, err := balancer.Invoke(ctx, fmt.Sprintf("request-%d", i))
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		fmt.Printf("Request %d: %v\n", i, resp)
	}
}

func testBalancer3() {
	addrs := []string{"backend1:8080", "backend2:8080", "backend3:8080"}
	balancer := NewSmartBalancer(addrs, 3)

	ctx := context.Background()

	for i := 0; i < 6; i++ {
		resp, err := balancer.Invoke(ctx, fmt.Sprintf("request-%d", i))
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		fmt.Printf("Request %d: %v\n", i, resp)
	}
}

func testConcurrency() {
	addrs := []string{"backend1:8080", "backend2:8080", "backend3:8080"}
	balancer := NewAtomicBalancer(addrs)

	var wg sync.WaitGroup
	ctx := context.Background()

	// Запускаем 10 горутин, каждая делает 3 запроса.
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 3; j++ {
				resp, err := balancer.Invoke(ctx, fmt.Sprintf("goroutine-%d-req-%d", id, j))
				if err != nil {
					fmt.Printf("Error in goroutine %d: %v\n", id, err)
					continue
				}
				fmt.Printf("Goroutine %d, Request %d: %v\n", id, j, resp)
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("Все горутины завершены")
}

// Объяснение решений:
//
// ============================================================================
// Решение 1: Простой Round Robin с sync.Mutex
// ============================================================================
//
// Алгоритм:
// 1. Храним индекс next - следующий backend для обработки
// 2. При каждом запросе берем backend[next]
// 3. Увеличиваем next по модулю количества backend'ов
// 4. Используем мьютекс для потокобезопасности
//
// Плюсы:
// - Простая и понятная реализация
// - Гарантированное равномерное распределение
// - Предсказуемое поведение
//
// Минусы:
// - Мьютекс может стать bottleneck при высокой нагрузке
// - Contention на мьютексе в конкурентной среде
//
// Сложность:
// - Время: O(1)
// - Память: O(n) где n = количество backend'ов
//
// ============================================================================
// Решение 2: Round Robin с atomic
// ============================================================================
//
// Алгоритм:
// 1. Используем atomic.AddUint64 для потокобезопасного счетчика
// 2. Каждая горутина получает уникальное значение счетчика
// 3. Вычисляем индекс backend'а: counter % len(backends)
// 4. Избегаем блокировок, используя lock-free операции
//
// Плюсы:
// - Значительно быстрее при высокой конкурентности
// - Нет contention на locks
// - Масштабируется лучше на многоядерных системах
//
// Минусы:
// - Чуть более сложная логика
// - Counter может переполниться (но это не проблема благодаря модулю)
//
// Производительность:
// - Atomic операции ~10-100x быстрее мьютекса при высокой конкурентности
// - Лучше использовать для production систем с высокой нагрузкой
//
// ============================================================================
// Решение 3: Smart Balancer с retry
// ============================================================================
//
// Алгоритм:
// 1. При ошибке автоматически пробуем следующий backend
// 2. Ограничиваем количество попыток (maxRetries)
// 3. Возвращаем успешный ответ или последнюю ошибку
//
// Плюсы:
// - Устойчивость к падениям отдельных backend'ов
// - Автоматический failover
// - Повышает availability системы
//
// Минусы:
// - Увеличивает latency при ошибках
// - Может маскировать проблемы с backend'ами
// - Нужна осторожность с idempotent операциями
//
// Когда использовать:
// - Backend'ы нестабильны
// - Важна высокая availability
// - Операции idempotent
//
// ============================================================================
// Сравнение алгоритмов балансировки
// ============================================================================
//
// Round Robin:
// - Простой и предсказуемый
// - Равномерное распределение
// - Не учитывает нагрузку на backend'ы
//
// Weighted Round Robin:
// - Учитывает мощность backend'ов
// - Распределяет пропорционально весам
// - Требует настройки весов
//
// Least Connections:
// - Направляет на наименее загруженный backend
// - Требует отслеживания активных соединений
// - Сложнее реализация
//
// Random:
// - Выбирает случайный backend
// - Простая реализация
// - Статистически равномерное распределение
//
// Consistent Hashing:
// - Стабильное распределение по ключу
// - Минимальная реорганизация при добавлении/удалении backend'ов
// - Используется в кешах и шардированных БД
//
// ============================================================================
// Дополнительные улучшения
// ============================================================================
//
// 1. Health Checks:
//    - Периодическая проверка доступности backend'ов
//    - Исключение недоступных из ротации
//    - Автоматическое восстановление при recovery
//
// 2. Circuit Breaker:
//    - Быстрый fail при недоступности backend'а
//    - Избегаем cascade failures
//    - Периодические попытки восстановления
//
// 3. Adaptive Load Balancing:
//    - Мониторинг latency и error rate
//    - Динамическое изменение весов
//    - Направление нагрузки на более быстрые backend'ы
//
// 4. Sticky Sessions:
//    - Привязка сессии к конкретному backend'у
//    - Для stateful приложений
//    - Consistent hashing по session ID
//
// 5. Metrics и Monitoring:
//    - Счетчики запросов на каждый backend
//    - Latency и error rate
//    - Alerting при проблемах
//
// ============================================================================
// Производственные рекомендации
// ============================================================================
//
// 1. Используйте AtomicBalancer для высоких нагрузок
// 2. Добавьте health checks для надежности
// 3. Реализуйте circuit breaker для защиты от cascade failures
// 4. Логируйте метрики для мониторинга
// 5. Используйте context для timeout и cancellation
// 6. Учитывайте idempotency при retry
// 7. Тестируйте под нагрузкой (load testing)
//
// ============================================================================
// Примеры использования в production
// ============================================================================
//
// - Nginx: Round Robin, Least Connections, IP Hash
// - HAProxy: Round Robin, Least Connections, Source IP
// - Kubernetes: Round Robin для Services
// - gRPC: Round Robin, Pick First, Custom policies
// - Envoy: Round Robin, Least Request, Ring Hash
