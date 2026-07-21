package main

import "fmt"

type WordCounter struct {
	counts map[string]int
	order  []string // Для отслеживания порядка добавления
	limit  int
}

func NewWordCounter(limit int) *WordCounter {
	return &WordCounter{
		counts: make(map[string]int),
		order:  make([]string, 0, limit),
		limit:  limit,
	}
}

func (wc *WordCounter) CountWord(word string) {
	// Если слово уже есть в мапе, просто увеличиваем счетчик.
	// Порядок добавления не изменяется - удаляем самые старые по времени добавления.
	if _, exists := wc.counts[word]; exists {
		wc.counts[word]++
		return
	}

	// Новое слово: добавляем в мапу и в порядок.
	wc.counts[word] = 1
	wc.order = append(wc.order, word)

	// Если превысили лимит, удаляем самое старое слово (первое в order).
	if len(wc.counts) > wc.limit {
		// Удаляем первое (самое старое) слово из order.
		oldest := wc.order[0]
		wc.order = wc.order[1:]
		delete(wc.counts, oldest)
	}
}

func main() {
	wc := NewWordCounter(3)

	words := []string{"apple", "banana", "apple", "orange", "grape", "banana", "kiwi"}
	for _, word := range words {
		wc.CountWord(word)
	}

	fmt.Println("Количество слов:", wc.counts)
}

// Ответ:
// Количество слов: map[grape:1 kiwi:1 orange:1]
//
// Пошаговое объяснение:
// 1. apple: добавляется (counts: {apple:1}, order: [apple], len=1)
// 2. banana: добавляется (counts: {apple:1, banana:1}, order: [apple, banana], len=2)
// 3. apple: увеличивается счетчик до 2, порядок не меняется (counts: {apple:2, banana:1}, order: [apple, banana], len=2)
// 4. orange: добавляется (counts: {apple:2, banana:1, orange:1}, order: [apple, banana, orange], len=3)
// 5. grape: добавляется, len=4 > limit=3 → удаляем самое старое (apple)
//    (counts: {banana:1, orange:1, grape:1}, order: [banana, orange, grape], len=3)
// 6. banana: увеличивается счетчик до 2, порядок не меняется (counts: {banana:2, orange:1, grape:1}, order: [banana, orange, grape], len=3)
// 7. kiwi: добавляется, len=4 > limit=3 → удаляем самое старое (banana)
//    (counts: {orange:1, grape:1, kiwi:1}, order: [orange, grape, kiwi], len=3)
//
// Итоговый результат: map[grape:1 kiwi:1 orange:1]
//
// Важно: удаляем самые старые записи по времени добавления, независимо от того,
// использовались ли они недавно или нет. Порядок добавления не обновляется при увеличении счетчика.

// Альтернативное решение с использованием LRU кэша (более сложное, но эффективное):
// Можно использовать структуру данных типа LRU cache для автоматического
// удаления наименее используемых элементов.
