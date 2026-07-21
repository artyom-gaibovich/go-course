package main

import (
	"fmt"
	"sync"
	"time"
)

// Service интерфейс для аналитики в реальном времени.
type Service interface {
	HandleEvent(userName string, currentTime time.Time)
	GetCount(userName string, currentTime time.Time) int
}

// ============================================================================
// Решение 1: Базовая реализация (очистка при чтении)
// ============================================================================

// analyticsService реализует Service.
// Хранит временные метки событий для каждого пользователя.
type analyticsService struct {
	mu     sync.RWMutex           // Защита от конкурентного доступа
	events map[string][]time.Time // userName -> список временных меток
	window time.Duration          // Временное окно (5 минут)
}

// NewAnalyticsService создает новый экземпляр сервиса аналитики.
func NewAnalyticsService() Service {
	return &analyticsService{
		events: make(map[string][]time.Time),
		window: 5 * time.Minute, // Временное окно 5 минут
	}
}

// HandleEvent регистрирует событие пользователя.
// Метод потокобезопасен и может вызываться из нескольких горутин.
func (s *analyticsService) HandleEvent(userName string, currentTime time.Time) {
	s.mu.Lock()         // Блокируем для записи
	defer s.mu.Unlock() // Разблокируем после выхода

	// Добавляем временную метку события.
	s.events[userName] = append(s.events[userName], currentTime)
}

// GetCount возвращает количество событий пользователя за последние 5 минут.
// Метод потокобезопасен и может вызываться из нескольких горутин.
// ВАЖНО: Очистка старых данных происходит здесь, при чтении.
func (s *analyticsService) GetCount(userName string, currentTime time.Time) int {
	s.mu.Lock()         // Блокируем для записи (так как будем очищать старые)
	defer s.mu.Unlock() // Разблокируем после выхода

	timestamps, exists := s.events[userName]
	if !exists {
		return 0 // Пользователь не найден
	}

	// Вычисляем границу временного окна.
	windowStart := currentTime.Add(-s.window)

	// Находим первый элемент, который попадает в окно.
	// Все элементы до него можно удалить.
	validStart := 0
	for i, ts := range timestamps {
		if ts.After(windowStart) || ts.Equal(windowStart) {
			validStart = i
			break
		}
	}

	// Удаляем устаревшие события для экономии памяти.
	// Создаем новый слайс с валидными событиями.
	validEvents := timestamps[validStart:]
	s.events[userName] = validEvents

	// Возвращаем количество событий в окне.
	return len(validEvents)
}

// ============================================================================
// Решение 2: С фоновой очисткой (background cleanup)
// ============================================================================

// analyticsServiceWithCleanup реализует Service с фоновой очисткой.
// Периодически удаляет старые события, снижая нагрузку на GetCount.
type analyticsServiceWithCleanup struct {
	mu     sync.RWMutex
	events map[string][]time.Time
	window time.Duration

	// Для фоновой очистки
	cleanupInterval time.Duration
	done            chan struct{}
	wg              sync.WaitGroup
}

// NewAnalyticsServiceWithCleanup создает сервис с фоновой очисткой.
func NewAnalyticsServiceWithCleanup(cleanupInterval time.Duration) *analyticsServiceWithCleanup {
	s := &analyticsServiceWithCleanup{
		events:          make(map[string][]time.Time),
		window:          5 * time.Minute,
		cleanupInterval: cleanupInterval,
		done:            make(chan struct{}),
	}

	// Запускаем фоновую горутину для очистки
	s.wg.Add(1)
	go s.backgroundCleanup()

	return s
}

func (s *analyticsServiceWithCleanup) HandleEvent(userName string, currentTime time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events[userName] = append(s.events[userName], currentTime)
}

func (s *analyticsServiceWithCleanup) GetCount(userName string, currentTime time.Time) int {
	s.mu.RLock() // Используем RLock, так как только читаем
	defer s.mu.RUnlock()

	timestamps, exists := s.events[userName]
	if !exists {
		return 0
	}

	windowStart := currentTime.Add(-s.window)

	// Считаем валидные события без модификации слайса
	count := 0
	for _, ts := range timestamps {
		if ts.After(windowStart) || ts.Equal(windowStart) {
			count++
		}
	}

	return count
}

// backgroundCleanup периодически очищает старые события.
func (s *analyticsServiceWithCleanup) backgroundCleanup() {
	defer s.wg.Done()

	ticker := time.NewTicker(s.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.cleanup()
		case <-s.done:
			return
		}
	}
}

// cleanup удаляет устаревшие события для всех пользователей.
func (s *analyticsServiceWithCleanup) cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-s.window)

	for userName, timestamps := range s.events {
		// Находим первый валидный элемент
		validStart := 0
		for i, ts := range timestamps {
			if ts.After(windowStart) || ts.Equal(windowStart) {
				validStart = i
				break
			}
		}

		// Обновляем слайс
		if validStart > 0 {
			s.events[userName] = timestamps[validStart:]
		}

		// Удаляем пользователя, если у него нет событий
		if len(s.events[userName]) == 0 {
			delete(s.events, userName)
		}
	}
}

// Stop останавливает фоновую горутину.
func (s *analyticsServiceWithCleanup) Stop() {
	close(s.done)
	s.wg.Wait()
}

func main() {
	fmt.Println("=== Решение 1: Базовая реализация ===")
	demonstrateBasic()

	fmt.Println("\n=== Решение 2: С фоновой очисткой ===")
	demonstrateWithCleanup()
}

func demonstrateBasic() {
	service := NewAnalyticsService()

	now := time.Now()

	// Регистрируем события
	service.HandleEvent("user1", now.Add(-6*time.Minute)) // Старое, не попадет
	service.HandleEvent("user1", now.Add(-4*time.Minute))
	service.HandleEvent("user1", now.Add(-3*time.Minute))
	service.HandleEvent("user1", now.Add(-1*time.Minute))
	service.HandleEvent("user2", now.Add(-2*time.Minute))

	// Получаем статистику
	fmt.Printf("user1 events: %d\n", service.GetCount("user1", now)) // 3
	fmt.Printf("user2 events: %d\n", service.GetCount("user2", now)) // 1
	fmt.Printf("user3 events: %d\n", service.GetCount("user3", now)) // 0
}

func demonstrateWithCleanup() {
	service := NewAnalyticsServiceWithCleanup(1 * time.Second)
	defer service.Stop()

	now := time.Now()

	// Регистрируем события
	service.HandleEvent("user1", now.Add(-6*time.Minute))
	service.HandleEvent("user1", now.Add(-4*time.Minute))
	service.HandleEvent("user1", now.Add(-3*time.Minute))
	service.HandleEvent("user1", now.Add(-1*time.Minute))
	service.HandleEvent("user2", now.Add(-2*time.Minute))

	// Даем время на фоновую очистку
	time.Sleep(1500 * time.Millisecond)

	// Получаем статистику
	fmt.Printf("user1 events: %d\n", service.GetCount("user1", now)) // 3
	fmt.Printf("user2 events: %d\n", service.GetCount("user2", now)) // 1
	fmt.Printf("user3 events: %d\n", service.GetCount("user3", now)) // 0
}

// Объяснение решений:
//
// ============================================================================
// Решение 1: Базовая реализация (cleanup on read)
// ============================================================================
//
// ПОДХОД: Очистка старых данных происходит в GetCount при чтении.
//
// Плюсы:
// - Простая реализация
// - Нет дополнительных горутин
// - Очистка происходит только для запрашиваемых пользователей
// - Lazy approach - не тратим ресурсы на неактивных пользователей
//
// Минусы:
// - GetCount выполняет двойную работу: подсчет + очистка
// - Используется Lock вместо RLock в GetCount
// - Память может расти между вызовами GetCount
// - Для редко запрашиваемых пользователей данные накапливаются
//
// Когда использовать:
// - Простые приложения
// - Небольшое количество пользователей
// - GetCount вызывается регулярно
// - Не критично использование Lock
//
// Сложность:
// - HandleEvent: O(1) append
// - GetCount: O(n) где n = количество событий пользователя
//
// ============================================================================
// Решение 2: С фоновой очисткой (background cleanup)
// ============================================================================
//
// ПОДХОД: Фоновая горутина периодически очищает старые данные для всех
//         пользователей. GetCount только читает и считает.
//
// Плюсы:
// - GetCount быстрее - только подсчет, без модификации
// - Используется RLock в GetCount - лучшая конкурентность
// - Проактивная очистка памяти
// - Предсказуемое использование ресурсов
//
// Минусы:
// - Дополнительная горутина (требует управления жизненным циклом)
// - Очистка всех пользователей каждый раз (даже неактивных)
// - Чуть более сложная реализация
// - Нужно корректно останавливать через Stop()
//
// Когда использовать:
// - Высокая нагрузка на GetCount
// - Много пользователей
// - Важна скорость чтения
// - Доступна память для дополнительной горутины
//
// Сложность:
// - HandleEvent: O(1) append
// - GetCount: O(n) где n = количество событий пользователя
// - Cleanup: O(N × M) где N = пользователей, M = событий на пользователя
//
// ============================================================================
// Сравнение подходов
// ============================================================================
//
// Метрика                | Базовая    | С фоновой очисткой
// -----------------------|------------|-------------------
// Простота реализации    | ✅ Высокая | ⚠️ Средняя
// Скорость GetCount      | ⚠️ Средняя | ✅ Высокая
// Использование памяти   | ⚠️ Растет  | ✅ Контролируемое
// Конкурентность чтения  | ❌ Lock    | ✅ RLock
// Дополнительные горутины| ✅ Нет     | ❌ Да (+1)
// Graceful shutdown      | ✅ Простой | ⚠️ Требует Stop()
//
// ============================================================================
// Настройка cleanupInterval
// ============================================================================
//
// Выбор интервала зависит от нагрузки и требований:
//
// cleanupInterval = 10 секунд:
//   - Редкая очистка
//   - Меньше нагрузки на Lock
//   - Больше накопленных данных между очистками
//   - Подходит для: низкая нагрузка, много памяти
//
// cleanupInterval = 1 секунда:
//   - Частая очистка
//   - Больше нагрузки на Lock
//   - Минимум накопленных данных
//   - Подходит для: высокая нагрузка, мало памяти
//
// Правило выбора:
//   cleanupInterval ≈ window / 10
//   Для window=5 минут → cleanupInterval=30 секунд
//
// ============================================================================
// Оптимизации
// ============================================================================
//
// 1. Sharding по пользователям:
//    Разделить пользователей по нескольким maps с отдельными мьютексами.
//    Снижает contention на мьютексе.
//
// 2. Бинарный поиск:
//    Если события отсортированы, использовать бинарный поиск границы окна.
//    Сложность: O(log n) вместо O(n)
//
// 3. Bucket approach:
//    Группировать события по временным bucket'ам (например, по секундам).
//    Хранить счетчики вместо всех timestamps.
//    Меньше памяти, но менее точно.
//
// 4. Adaptive cleanup:
//    Динамически изменять cleanupInterval в зависимости от нагрузки.
//    Высокая нагрузка → чаще очистка.
//
// ============================================================================
// КЛЮЧЕВЫЕ ВЫВОДЫ
// ============================================================================
//
// - Базовая реализация: простая, для небольших нагрузок
// - Фоновая очистка: для высоких нагрузок и множества пользователей
// - Trade-off: простота vs производительность
// - Выбор зависит от требований: RPS, количество пользователей, память
// - Фоновая очистка требует graceful shutdown через Stop()
