package main

import (
	"fmt"
	"sync"
)

// Document представляет документ с версионированием
type Document struct {
	Url            string
	PubDate        uint64
	FetchTime      uint64
	Text           string
	FirstFetchTime *uint64
}

type Processor interface {
	Process(doc Document) (*Document, error)
}

// ============================================================================
// Решение: Document Versioning Processor
// ============================================================================

// documentState хранит состояние документа по URL
type documentState struct {
	url             string
	maxFetchTime    uint64              // Максимальный FetchTime (последняя версия)
	minFetchTime    uint64              // Минимальный FetchTime (первая версия)
	latestText      string              // Текст из версии с максимальным FetchTime
	earliestPubDate uint64              // PubDate из версии с минимальным FetchTime
	seenFetchTimes  map[uint64]struct{} // Для отслеживания дубликатов
}

type documentProcessor struct {
	mu   sync.RWMutex
	docs map[string]*documentState
}

// NewDocumentProcessor создает новый процессор документов
func NewDocumentProcessor() Processor {
	return &documentProcessor{
		docs: make(map[string]*documentState),
	}
}

// Process обрабатывает входящий документ
func (p *documentProcessor) Process(doc Document) (*Document, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Получаем или создаем состояние для данного URL
	state, exists := p.docs[doc.Url]
	if !exists {
		// Первый документ для этого URL
		state = &documentState{
			url:             doc.Url,
			maxFetchTime:    doc.FetchTime,
			minFetchTime:    doc.FetchTime,
			latestText:      doc.Text,
			earliestPubDate: doc.PubDate,
			seenFetchTimes:  make(map[uint64]struct{}),
		}
		state.seenFetchTimes[doc.FetchTime] = struct{}{}
		p.docs[doc.Url] = state

		// Формируем результат
		firstFetchTime := doc.FetchTime
		return &Document{
			Url:            doc.Url,
			PubDate:        doc.PubDate,
			FetchTime:      doc.FetchTime,
			Text:           doc.Text,
			FirstFetchTime: &firstFetchTime,
		}, nil
	}

	// Проверяем, не дубликат ли это
	if _, ok := state.seenFetchTimes[doc.FetchTime]; ok {
		return nil, nil // Дубликат, не требует обновления
	}
	state.seenFetchTimes[doc.FetchTime] = struct{}{}

	// Если это более свежая версия (больший FetchTime)
	if doc.FetchTime > state.maxFetchTime {
		state.maxFetchTime = doc.FetchTime
		state.latestText = doc.Text
	}

	// Если это более старая версия (меньший FetchTime)
	if doc.FetchTime < state.minFetchTime {
		state.minFetchTime = doc.FetchTime
		state.earliestPubDate = doc.PubDate
	}

	return &Document{
		Url:            state.url,
		PubDate:        state.earliestPubDate,
		FetchTime:      state.maxFetchTime,
		Text:           state.latestText,
		FirstFetchTime: &state.minFetchTime,
	}, nil
}

// ============================================================================
// Демонстрация и тестирование
// ============================================================================

func main() {
	fmt.Println("=== Тест версионирования документов ===\n")

	processor := NewDocumentProcessor()

	// Тест 1: Первый документ
	fmt.Println("Тест 1: Первый документ для URL")
	doc1 := Document{
		Url:       "doc1",
		FetchTime: 100,
		Text:      "Version 1",
		PubDate:   50,
	}
	result := processAndPrint(processor, doc1)
	fmt.Printf("Результат: Text=%s, FetchTime=%d, PubDate=%d, FirstFetchTime=%d\n",
		result.Text, result.FetchTime, result.PubDate, *result.FirstFetchTime)
	fmt.Println()

	// Тест 2: Более новая версия (больший FetchTime)
	fmt.Println("Тест 2: Более новая версия")
	doc2 := Document{
		Url:       "doc1",
		FetchTime: 200,
		Text:      "Version 2",
		PubDate:   60,
	}
	result = processAndPrint(processor, doc2)
	if result != nil {
		fmt.Printf("Результат: Text=%s, FetchTime=%d, PubDate=%d, FirstFetchTime=%d\n",
			result.Text, result.FetchTime, result.PubDate, *result.FirstFetchTime)
		fmt.Println("✓ Текст обновился на 'Version 2', FetchTime=200, PubDate остался 50")
	}
	fmt.Println()

	// Тест 3: Более старая версия (меньший FetchTime)
	fmt.Println("Тест 3: Более старая версия")
	doc3 := Document{
		Url:       "doc1",
		FetchTime: 50,
		Text:      "Version 0",
		PubDate:   40,
	}
	result = processAndPrint(processor, doc3)
	if result != nil {
		fmt.Printf("Результат: Text=%s, FetchTime=%d, PubDate=%d, FirstFetchTime=%d\n",
			result.Text, result.FetchTime, result.PubDate, *result.FirstFetchTime)
		fmt.Println("✓ Текст остался 'Version 2', но PubDate=40 и FirstFetchTime=50")
	}
	fmt.Println()

	// Тест 4: Дубликат
	fmt.Println("Тест 4: Дубликат документа")
	doc4 := Document{
		Url:       "doc1",
		FetchTime: 200, // Уже был
		Text:      "Version 2 duplicate",
		PubDate:   70,
	}
	result = processAndPrint(processor, doc4)
	if result == nil {
		fmt.Println("✓ Дубликат корректно отфильтрован")
	}
	fmt.Println()

	// Тест 5: Конкурентная обработка
	fmt.Println("Тест 5: Конкурентная обработка")
	testConcurrency()
}

func processAndPrint(processor Processor, doc Document) *Document {
	result, err := processor.Process(doc)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return nil
	}
	if result == nil {
		fmt.Println("Результат: nil (дубликат или нет изменений)")
		return nil
	}
	return result
}

func testConcurrency() {
	processor := NewDocumentProcessor()
	var wg sync.WaitGroup

	// Запускаем 10 горутин, каждая отправляет документы
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for j := 0; j < 100; j++ {
				doc := Document{
					Url:       fmt.Sprintf("doc_%d", id),
					FetchTime: uint64(j * 10),
					Text:      fmt.Sprintf("Text from goroutine %d, iteration %d", id, j),
					PubDate:   uint64(j * 5),
				}
				processor.Process(doc)
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("✓ Конкурентная обработка завершена без ошибок")
}

// Объяснение решения:
//
// ============================================================================
// Ключевые концепции
// ============================================================================
//
// 1. ВЕРСИОНИРОВАНИЕ ДОКУМЕНТОВ
//
// Проблема: Документы приходят из распределенной системы (Kafka) в произвольном
// порядке. Нужно собрать консистентное представление документа.
//
// Решение: Храним состояние документа и отслеживаем:
// - maxFetchTime: для получения последнего текста
// - minFetchTime: для определения первой версии
// - earliestPubDate: для сохранения исходной даты публикации
// - seenFetchTimes: для фильтрации дубликатов
//
// 2. ОБРАБОТКА РАЗЛИЧНЫХ СЦЕНАРИЕВ
//
// Сценарий A: Документы приходят по порядку (FetchTime: 100, 200, 300)
//   Каждый новый документ обновляет maxFetchTime и Text
//   minFetchTime и PubDate остаются от первого документа
//
// Сценарий B: Документы приходят не по порядку (200, 100, 300)
//   При FetchTime=200: устанавливаем max и min
//   При FetchTime=100: обновляем min и PubDate, max остается 200
//   При FetchTime=300: обновляем max и Text
//
// Сценарий C: Дубликаты (100, 200, 100, 200)
//   Отслеживаем через seenFetchTimes map
//   Повторные FetchTime игнорируются, возвращаем nil
//
// 3. КОНКУРЕНТНЫЙ ДОСТУП
//
// Проблема: Метод Process вызывается из нескольких горутин одновременно
//
// Решение: sync.RWMutex
// - Lock() для модификации состояния
// - Защита map и documentState
//
// Альтернативы:
// - sync.Map: если много разных URL, мало конфликтов
// - Channel-based: для строгой последовательности обработки
// - Sharding: разделение по URL для снижения contention
//
// 4. ОТСЛЕЖИВАНИЕ ДУБЛИКАТОВ
//
// Подход: map[uint64]bool для seenFetchTimes
//
// Плюсы:
// - O(1) проверка и вставка
// - Простая реализация
//
// Минусы:
// - Память растет с количеством версий
// - Для долгоживущих документов может быть проблемой
//
// Оптимизации:
// - Ограничить размер map (удалять старые FetchTime)
// - Использовать Bloom Filter для вероятностной проверки
// - TTL для автоочистки старых записей
//
// 5. ВОЗВРАТ РЕЗУЛЬТАТА
//
// nil возвращается когда:
// - Дубликат (FetchTime уже был)
// - Документ не вносит изменений (промежуточный FetchTime без нового PubDate)
//
// Document возвращается когда:
// - Новый URL (первый документ)
// - Обновился maxFetchTime (новый Text)
// - Обновился minFetchTime (новый FirstFetchTime или PubDate)
//
// ============================================================================
// Сложность
// ============================================================================
//
// Время:
// - Process: O(1) амортизированное
// - Map lookup/insert: O(1)
// - Сравнения и обновления: O(1)
//
// Память:
// - O(U × V) где U = количество уникальных URL, V = версий на URL
// - seenFetchTimes растет с каждой новой версией
//
// ============================================================================
// Оптимизации для production
// ============================================================================
//
// 1. SHARDING ПО URL
//
// Разделить документы по нескольким процессорам:
//
// type ShardedProcessor struct {
//     shards    []*documentProcessor
//     shardCount int
// }
//
// func (sp *ShardedProcessor) Process(doc Document) (*Document, error) {
//     shardIdx := hash(doc.Url) % sp.shardCount
//     return sp.shards[shardIdx].Process(doc)
// }
//
// Плюсы: снижает contention на мьютексе
//
// 2. ОГРАНИЧЕНИЕ ПАМЯТИ
//
// Лимит на количество хранимых FetchTime:
//
// const maxVersionsPerDoc = 1000
//
// if len(state.seenFetchTimes) > maxVersionsPerDoc {
//     // Удаляем старейшие FetchTime
//     deleteOldestFetchTimes(state)
// }
//
// 3. ПЕРИОДИЧЕСКАЯ ОЧИСТКА
//
// Фоновая горутина удаляет старые документы:
//
// go func() {
//     ticker := time.NewTicker(1 * time.Hour)
//     for range ticker.C {
//         p.cleanup(time.Now().Add(-24 * time.Hour))
//     }
// }()
//
// 4. МЕТРИКИ И МОНИТОРИНГ
//
// - Количество уникальных URL
// - Средний размер seenFetchTimes
// - Процент дубликатов
// - Latency обработки
//
// var (
//     documentsProcessed = prometheus.NewCounter(...)
//     duplicatesFiltered = prometheus.NewCounter(...)
//     processingDuration = prometheus.NewHistogram(...)
// )
//
// ============================================================================
// Альтернативные подходы
// ============================================================================
//
// ПОДХОД 1: Event Sourcing
//
// Храним все события (версии документов) в порядке получения
// При запросе - пересчитываем состояние
//
// Плюсы:
// - Полная история изменений
// - Возможность replay
//
// Минусы:
// - Медленные запросы (нужно пересчитывать)
// - Больше памяти
//
// ПОДХОД 2: CQRS (Command Query Responsibility Segregation)
//
// Разделить запись и чтение:
// - Write: сохраняем все версии
// - Read: отдельная структура с агрегированным состоянием
//
// Плюсы:
// - Оптимизировано под запись и чтение отдельно
// - Масштабируемость
//
// Минусы:
// - Eventual consistency
// - Сложнее в реализации
//
// ПОДХОД 3: Materialized View
//
// Храним финальное состояние в БД/кеше
// Process только обновляет это состояние
//
// Плюсы:
// - Персистентность
// - Быстрые запросы
//
// Минусы:
// - Дополнительная инфраструктура (БД)
// - Network latency
//
// ============================================================================
// Ключевые выводы
// ============================================================================
//
// 1. Храним состояние документа, отслеживая min/max FetchTime
// 2. Используем map для фильтрации дубликатов
// 3. Защищаем конкурентный доступ через sync.RWMutex
// 4. Возвращаем nil для дубликатов и неизменных документов
// 5. Для production: sharding, ограничение памяти, мониторинг
