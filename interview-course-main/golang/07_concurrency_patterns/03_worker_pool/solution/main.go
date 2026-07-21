package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// WorkerPool - пул воркеров для обработки задач.
type WorkerPool struct {
	workers     int
	tasks       chan string
	wg          sync.WaitGroup
	processFunc func(string) error
	closed      atomic.Bool // Атомарный флаг закрытия пула (Go 1.19+)
	closeOnce   sync.Once
}

// NewWorkerPool создает новый пул воркеров.
func NewWorkerPool(workers int, processFunc func(string) error) *WorkerPool {
	return &WorkerPool{
		workers:     workers,
		tasks:       make(chan string),
		processFunc: processFunc,
	}
}

// Start запускает воркеров.
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// worker обрабатывает задачи из канала.
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	for task := range wp.tasks {
		fmt.Printf("Worker %d обрабатывает: %s\n", id, task)
		if err := wp.processFunc(task); err != nil {
			fmt.Printf("Worker %d ошибка при обработке %s: %v\n", id, task, err)
		}
	}
}

// Submit отправляет задачу в пул.
// Возвращает ошибку, если пул закрыт.
func (wp *WorkerPool) Submit(task string) error {
	// Проверяем, не закрыт ли пул.
	if wp.closed.Load() {
		return fmt.Errorf("пул закрыт")
	}

	// Пытаемся отправить задачу с защитой от паники при закрытом канале.
	// Используем recover для перехвата паники от отправки в закрытый канал.
	var panicErr interface{}
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicErr = r
			}
		}()
		wp.tasks <- task
	}()

	// Если произошла паника, значит канал закрыт.
	if panicErr != nil {
		return fmt.Errorf("пул закрыт: %v", panicErr)
	}

	return nil
}

// Close закрывает пул и ждет завершения всех воркеров.
// Гарантирует, что канал будет закрыт только один раз.
func (wp *WorkerPool) Close() {
	wp.closeOnce.Do(func() {
		// Устанавливаем флаг закрытия.
		wp.closed.Store(true)
		// Закрываем канал - это остановит все воркеры.
		close(wp.tasks)
		// Ждем завершения всех воркеров.
		wp.wg.Wait()
	})
}

// processImage симулирует обработку изображения.
func processImage(imagePath string) error {
	// Симулируем дорогостоящую операцию
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("Изображение %s обработано\n", imagePath)
	return nil
}

func main() {
	// Создаем пул из 3 воркеров
	pool := NewWorkerPool(3, processImage)
	pool.Start()

	// Отправляем задачи
	images := []string{"img1.jpg", "img2.jpg", "img3.jpg", "img4.jpg", "img5.jpg"}
	for _, img := range images {
		if err := pool.Submit(img); err != nil {
			fmt.Printf("Ошибка отправки задачи %s: %v\n", img, err)
		}
	}

	// Закрываем пул и ждем завершения
	pool.Close()
	fmt.Println("Все изображения обработаны")
}

// Объяснение:
// 1. Worker Pool ограничивает количество одновременно работающих горутин.
// 2. Задачи отправляются в канал tasks.
// 3. Воркеры читают задачи из канала и обрабатывают их.
// 4. Количество воркеров фиксировано, что предотвращает перегрузку системы.
// 5. После отправки всех задач закрываем канал и ждем завершения воркеров.
//
// Защита от паники:
// - Используем atomic.Bool для атомарной проверки состояния пула (Go 1.19+).
//   Альтернатива для старых версий: int32 с atomic.LoadInt32/StoreInt32.
// - sync.Once гарантирует, что канал будет закрыт только один раз.
// - Submit проверяет флаг закрытия перед отправкой задачи.
// - recover перехватывает панику от отправки в закрытый канал.
