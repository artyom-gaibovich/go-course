package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

// ============================================================================
// Типы данных
// ============================================================================

type Document struct {
	DocID   string
	Content string
}

type User struct {
	UserID string
}

type Scorer interface {
	GetScore(doc Document, user User) int
}

type DocumentService interface {
	AddDocument(doc Document) error
	GetTopDocuments(user User, limit int) ([]Document, error)
}

// ============================================================================
// Решение 1: Naive (пересчет при каждом запросе)
// ============================================================================

type naiveDocumentService struct {
	mu     sync.RWMutex
	docs   map[string]Document // docID -> document
	scorer Scorer
}

func NewNaiveDocumentService(scorer Scorer) DocumentService {
	return &naiveDocumentService{
		docs:   make(map[string]Document),
		scorer: scorer,
	}
}

func (s *naiveDocumentService) AddDocument(doc Document) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.docs[doc.DocID] = doc
	return nil
}

func (s *naiveDocumentService) GetTopDocuments(user User, limit int) ([]Document, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Вычисляем скор для всех документов
	type scoredDoc struct {
		doc   Document
		score int
	}

	scored := make([]scoredDoc, 0, len(s.docs))
	for _, doc := range s.docs {
		score := s.scorer.GetScore(doc, user)
		scored = append(scored, scoredDoc{doc: doc, score: score})
	}

	// Сортируем по убыванию скора
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	// Берем топ K
	result := make([]Document, 0, limit)
	for i := 0; i < limit && i < len(scored); i++ {
		result = append(result, scored[i].doc)
	}

	return result, nil
}

// ============================================================================
// Решение 2: Precomputed Top-K с кешированием
// ============================================================================

type cachedDocumentService struct {
	mu     sync.RWMutex
	docs   map[string]Document   // docID -> document
	cache  map[string][]Document // userID -> top documents
	scorer Scorer
	topK   int // размер кеша
}

func NewCachedDocumentService(scorer Scorer, topK int) DocumentService {
	return &cachedDocumentService{
		docs:   make(map[string]Document),
		cache:  make(map[string][]Document),
		scorer: scorer,
		topK:   topK,
	}
}

func (s *cachedDocumentService) AddDocument(doc Document) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.docs[doc.DocID] = doc

	// Инвалидируем весь кеш при добавлении/обновлении документа.
	// Причина: новый/обновленный документ может изменить топ-K для любого пользователя.
	// Простое решение: очистить весь кеш и пересчитать при следующем запросе.
	//
	// Альтернатива (более сложная, но эффективная):
	// - Вычислить новый скор для обновленного документа для каждого пользователя
	// - Обновить кеш только если документ попал в топ-K или выпал из него
	// - Это требует O(U) операций, где U = количество пользователей в кеше
	//
	// Текущий подход проще и безопаснее: гарантирует актуальность данных.
	s.cache = make(map[string][]Document)

	return nil
}

func (s *cachedDocumentService) GetTopDocuments(user User, limit int) ([]Document, error) {
	s.mu.RLock()
	// Проверяем кеш
	if cached, exists := s.cache[user.UserID]; exists {
		s.mu.RUnlock()

		// Возвращаем из кеша с учетом limit
		if limit > len(cached) {
			return cached, nil
		}
		return cached[:limit], nil
	}
	s.mu.RUnlock()

	// Кеша нет, вычисляем
	s.mu.Lock()
	defer s.mu.Unlock()

	// Double-check (могло быть вычислено другой горутиной)
	if cached, exists := s.cache[user.UserID]; exists {
		if limit > len(cached) {
			return cached, nil
		}
		return cached[:limit], nil
	}

	// Вычисляем топ документов
	type scoredDoc struct {
		doc   Document
		score int
	}

	scored := make([]scoredDoc, 0, len(s.docs))
	for _, doc := range s.docs {
		score := s.scorer.GetScore(doc, user)
		scored = append(scored, scoredDoc{doc: doc, score: score})
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	// Сохраняем в кеш топ K
	topDocs := make([]Document, 0, s.topK)
	for i := 0; i < s.topK && i < len(scored); i++ {
		topDocs = append(topDocs, scored[i].doc)
	}
	s.cache[user.UserID] = topDocs

	// Возвращаем с учетом limit
	if limit > len(topDocs) {
		return topDocs, nil
	}
	return topDocs[:limit], nil
}

// ============================================================================
// Mock Scorer для тестирования
// ============================================================================

type mockScorer struct{}

func (s *mockScorer) GetScore(doc Document, user User) int {
	// Простой скоринг: количество слов в Content + длина UserID
	words := len(strings.Fields(doc.Content))
	userBonus := len(user.UserID) * 10
	return words + userBonus
}

// ============================================================================
// Демонстрация
// ============================================================================

func main() {
	scorer := &mockScorer{}

	fmt.Println("=== Решение 1: Naive ===")
	testService(NewNaiveDocumentService(scorer))

	fmt.Println("\n=== Решение 2: Cached ===")
	testService(NewCachedDocumentService(scorer, 100))
}

func testService(service DocumentService) {
	// Добавляем документы
	docs := []Document{
		{DocID: "1", Content: "Go programming language"},
		{DocID: "2", Content: "Python programming"},
		{DocID: "3", Content: "Java development guide"},
		{DocID: "4", Content: "Rust systems programming language tutorial"},
		{DocID: "5", Content: "JavaScript web development"},
	}

	for _, doc := range docs {
		service.AddDocument(doc)
	}

	// Получаем топ документов
	user := User{UserID: "user123"}
	top, err := service.GetTopDocuments(user, 3)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	fmt.Println("Топ 3 документа:")
	for i, doc := range top {
		fmt.Printf("%d. DocID=%s, Content=%s\n", i+1, doc.DocID, doc.Content)
	}
}

// Объяснение решений:
//
// ============================================================================
// Решение 1: Naive (пересчет при каждом запросе)
// ============================================================================
//
// Алгоритм:
// 1. Храним все документы в map
// 2. При GetTopDocuments:
//    - Вычисляем скор для ВСЕХ документов
//    - Сортируем по убыванию
//    - Берем топ K
//
// Плюсы:
// - Простая реализация
// - Всегда актуальные данные
// - Минимальное использование памяти
// - AddDocument O(1)
//
// Минусы:
// - GetTopDocuments O(N log N) где N = количество документов
// - Неэффективно при большом количестве документов
// - Вычисляем скор для всех, даже если нужен топ 10
//
// Когда использовать:
// - Малое количество документов (< 10K)
// - Редкие запросы GetTopDocuments
// - Частые обновления документов
//
// Сложность:
// - AddDocument: O(1)
// - GetTopDocuments: O(N log N)
// - Память: O(N)
//
// ============================================================================
// Решение 2: Precomputed Top-K с кешированием
// ============================================================================
//
// Алгоритм:
// 1. При первом GetTopDocuments для пользователя:
//    - Вычисляем топ K документов
//    - Сохраняем в кеш (map[userID][]Document)
// 2. Последующие запросы:
//    - Берем из кеша O(1)
// 3. При AddDocument:
//    - Инвалидируем весь кеш
//
// Плюсы:
// - Быстрый GetTopDocuments (из кеша) O(1)
// - Хорошо для read-heavy workload
// - Можно контролировать размер кеша (топ K)
//
// Минусы:
// - Инвалидация всего кеша при обновлении
// - Память растет с количеством пользователей O(U × K)
// - First request медленный (cold cache)
//
// Когда использовать:
// - Много пользователей, мало обновлений
// - Read-heavy workload (10:1 read:write)
// - Можем позволить задержку на первом запросе
//
// Сложность:
// - AddDocument: O(1) + инвалидация кеша
// - GetTopDocuments (cached): O(1)
// - GetTopDocuments (cold): O(N log N)
// - Память: O(N + U × K)
//
// Оптимизации:
// - Partial invalidation: инвалидировать только для пользователей,
//   для которых скор изменился (требует пересчета скора для всех пользователей)
// - TTL для кеша: автоматическая очистка старых записей
// - LRU eviction: удалять редко используемые записи
// - Lazy invalidation: не очищать кеш сразу, а помечать как "грязный"
//   и пересчитывать при следующем запросе
//
// Почему инвалидируем весь кеш:
// - Новый документ может попасть в топ-K для любого пользователя
// - Обновленный документ может изменить свой скор и позицию в топе
// - Простое решение гарантирует актуальность без сложной логики
// - Для read-heavy workload это приемлемо (редкие обновления)
//
// ============================================================================
// Сравнение решений
// ============================================================================
//
// Метрика            | Naive      | Cached     |
// -------------------|------------|------------|
// AddDocument        | O(1)       | O(1) + inv |
// GetTopDocuments    | O(N log N) | O(1)/O(N)  |
// Память             | O(N)       | O(N+U×K)   |
// Актуальность       | ✅ Всегда  | ⚠️ Кеш     |
// Read-heavy         | ❌ Плохо   | ✅ Отлично |
// Write-heavy        | ✅ Хорошо  | ❌ Плохо   |
//
// ============================================================================
// Реальные системы
// ============================================================================
//
// 1. ELASTICSEARCH
// - Использует inverted index + BM25 scoring
// - Top-K через aggregations
// - Sharding для масштабирования
//
// 2. REDIS SORTED SETS
// - ZADD для добавления с score
// - ZREVRANGE для топ K (O(log N + K))
// - Memory efficient
//
// 3. APPROXIMATE ALGORITHMS
// - Twitter: Summingbird для approximate top-K
// - Google: HyperLogLog для unique counts
// - Facebook: Count-Min Sketch
//
// ============================================================================
// Ключевые выводы
// ============================================================================
//
// 1. Trade-off между точностью и производительностью
// 2. Naive: простое, но медленное для больших N
// 3. Cached: быстрое чтение, медленная запись
// 4. Min-Heap: оптимальная сложность O(N log K)
// 5. Для production: inverted index или Redis Sorted Sets
// 6. Approximate алгоритмы для экстремальных масштабов
