package main

import (
	"fmt"
	"sync"
	"time"
)

// Semaphore - семафор для ограничения количества одновременных операций.
type Semaphore struct {
	ch chan struct{}
}

// NewSemaphore создает новый семафор с указанным лимитом.
func NewSemaphore(limit int) *Semaphore {
	return &Semaphore{
		ch: make(chan struct{}, limit),
	}
}

// Acquire занимает слот семафора.
func (s *Semaphore) Acquire() {
	s.ch <- struct{}{}
}

// Release освобождает слот семафора.
func (s *Semaphore) Release() {
	<-s.ch
}

// DownloadFile симулирует загрузку файла.
func DownloadFile(filename string, sem *Semaphore, wg *sync.WaitGroup) {
	defer wg.Done()

	// Занимаем слот семафора
	sem.Acquire()
	defer sem.Release()

	fmt.Printf("Начинаем загрузку %s\n", filename)
	// Симулируем загрузку
	time.Sleep(1 * time.Second)
	fmt.Printf("Загрузка %s завершена\n", filename)
}

func main() {
	// Создаем семафор с лимитом 3 одновременных соединений
	sem := NewSemaphore(3)
	var wg sync.WaitGroup

	files := []string{"file1.txt", "file2.txt", "file3.txt", "file4.txt", "file5.txt", "file6.txt"}

	for _, file := range files {
		wg.Add(1)
		go DownloadFile(file, sem, &wg)
	}

	wg.Wait()
	fmt.Println("Все файлы загружены")
}

// Объяснение:
// 1. Семафор ограничивает количество одновременных операций.
// 2. Acquire занимает слот (блокируется, если все слоты заняты).
// 3. Release освобождает слот.
// 4. Позволяет динамически контролировать количество активных операций.
// 5. Более гибкий, чем фиксированный пул воркеров.
