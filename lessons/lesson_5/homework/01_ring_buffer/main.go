/*
Домашнее задание: Кольцевой буфер (Ring Buffer)

Реализуйте обобщённую структуру RingBuffer[T any] — кольцевую очередь
фиксированного размера на основе слайса. Это FIFO-очередь, которая решает
проблему «медленного pop front» обычного слайса: все операции — O(1).

Интерфейс:

    type RingBuffer[T any] struct { ... }

    func NewRingBuffer[T any](capacity int) *RingBuffer[T]

    func (rb *RingBuffer[T]) Push(v T) error
        Добавить элемент в конец буфера.
        Вернуть ошибку (ErrFull), если буфер заполнен.

    func (rb *RingBuffer[T]) Pop() (T, error)
        Извлечь элемент из начала буфера.
        Вернуть ошибку (ErrEmpty), если буфер пуст.

    func (rb *RingBuffer[T]) Peek() (T, error)
        Прочитать элемент из начала буфера БЕЗ извлечения.
        Вернуть ошибку (ErrEmpty), если буфер пуст.

    func (rb *RingBuffer[T]) Len() int
    func (rb *RingBuffer[T]) Cap() int
    func (rb *RingBuffer[T]) IsFull() bool
    func (rb *RingBuffer[T]) IsEmpty() bool

Алгоритм:
  - Используйте два индекса: head (начало) и tail (следующая позиция для записи)
  - Backing array фиксированного размера создаётся один раз в NewRingBuffer
  - При достижении конца массива индексы «оборачиваются»: idx = (idx+1) % capacity
  - Слайс не должен расти через append — только preallocated backing array

Ограничения:
  - НЕ используйте append (буфер фиксированного размера)
  - Используйте дженерики (Go 1.21+)

Критерии:
  [ ] Все операции работают за O(1)
  [ ] Push возвращает ErrFull при полном буфере
  [ ] Pop возвращает ErrEmpty при пустом буфере
  [ ] Корректная работа после «оборачивания» индексов
  [ ] В main: демонстрация push/pop, переполнения и оборачивания

Пример поведения (capacity=3):
  Push(1) Push(2) Push(3) → [1,2,3], full=true
  Push(4)                 → error: buffer is full
  Pop()                   → 1, [2,3]
  Push(4)                 → [2,3,4]
  Pop(), Pop(), Pop()     → 2, 3, 4, empty=true
  Pop()                   → error: buffer is empty
*/

package main

import (
	"errors"
	"fmt"
)

var (
	ErrFull  = errors.New("buffer is full")
	ErrEmpty = errors.New("buffer is empty")
)

type RingBuffer[T any] struct {
	// TODO: добавьте поля
	data []T
	head int
	tail int
	size int
}

func NewRingBuffer[T any](capacity int) *RingBuffer[T] {
	// TODO: реализуйте
	return &RingBuffer[T]{data: make([]T, capacity)}
}

func (rb *RingBuffer[T]) Push(v T) error {
	// TODO: реализуйте
	if rb.IsFull() {
		return ErrFull
	}
	rb.data[rb.tail] = v
	rb.tail = (rb.tail + 1) % rb.Cap()
	rb.size++
	return nil
}

func (rb *RingBuffer[T]) Pop() (T, error) {
	// TODO: реализуйте
	var zero T
	if rb.IsEmpty() {
		return zero, ErrEmpty
	}
	value := rb.data[rb.head]
	rb.head = (rb.head + 1) % rb.Cap()
	rb.size--
	return value, nil
}

func (rb *RingBuffer[T]) Peek() (T, error) {
	// TODO: реализуйте
	var zero T
	if rb.IsEmpty() {
		return zero, ErrEmpty
	}
	return rb.data[rb.head], nil
}

func (rb *RingBuffer[T]) Len() int {
	// TODO: реализуйте
	return rb.size
}

func (rb *RingBuffer[T]) Cap() int {
	// TODO: реализуйте
	return len(rb.data)
}

func (rb *RingBuffer[T]) IsFull() bool {
	return rb.Len() == rb.Cap()
}

func (rb *RingBuffer[T]) IsEmpty() bool {
	return rb.Len() == 0
}

func main() {
	fmt.Println("=== RingBuffer[int], capacity=3 ===")
	rb := NewRingBuffer[int](3)

	fmt.Println("Push 1, 2, 3:")
	for _, v := range []int{1, 2, 3} {
		if err := rb.Push(v); err != nil {
			fmt.Println("  Push error:", err)
		} else {
			fmt.Printf("  pushed %d, len=%d, full=%v\n", v, rb.Len(), rb.IsFull())
		}
	}

	fmt.Println("Push 4 (должна быть ошибка):")
	if err := rb.Push(4); err != nil {
		fmt.Println("  ожидаемая ошибка:", err)
	}

	fmt.Println("Pop:")
	for !rb.IsEmpty() {
		v, err := rb.Pop()
		if err != nil {
			fmt.Println("  Pop error:", err)
		} else {
			fmt.Printf("  popped %d, len=%d\n", v, rb.Len())
		}
	}

	fmt.Println("Pop из пустого (должна быть ошибка):")
	if _, err := rb.Pop(); err != nil {
		fmt.Println("  ожидаемая ошибка:", err)
	}

	fmt.Println()
	fmt.Println("=== Оборачивание индексов: capacity=4 ===")
	rb2 := NewRingBuffer[string](4)
	words := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	for i, w := range words {
		if i%2 == 0 && rb2.Len() > 0 {
			v, _ := rb2.Pop()
			fmt.Printf("pop: %s | ", v)
		}
		if err := rb2.Push(w); err != nil {
			fmt.Printf("push(%s) err: %v | ", w, err)
		} else {
			fmt.Printf("push(%s) len=%d | ", w, rb2.Len())
		}
	}
	fmt.Println()
	fmt.Println("финальный drain:")
	for !rb2.IsEmpty() {
		v, _ := rb2.Pop()
		fmt.Printf("  %s\n", v)
	}
}
