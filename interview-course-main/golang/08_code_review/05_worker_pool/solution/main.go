package main

// СПИСОК НАЙДЕННЫХ ПРОБЛЕМ И ИСПРАВЛЕНИЙ:
//
// 1. Создается столько воркеров, сколько задач (строка 42)
//    Проблема: for range tasks { go worker() }
//              Если задач 1000 - создастся 1000 горутин!
//              Нет ограничения на количество параллельных воркеров
//    Решение: Создать фиксированное количество воркеров (например, NumCPU)
//
// 2. Канал jobChan никогда не закрывается (строка 47-51)
//    Проблема: go func() { for _, task := range tasks { jobChan <- task } }()
//              Канал не закрывается после отправки всех задач
//              Воркеры ждут в for job := range jobs бесконечно
//    Решение: Закрыть jobChan после отправки всех задач
//
// 3. Goroutine leak - воркеры не завершаются (строка 67-75)
//    Проблема: worker() выполняет for job := range jobs
//              Но jobChan не закрывается, воркеры висят навсегда
//    Решение: Закрыть jobChan, тогда range завершится
//
// 4. Каналы не буферизированы (строка 38-39)
//    Проблема: jobChan := make(chan Task)
//              Отправка блокируется до чтения
//              Может привести к deadlock если порядок неправильный
//    Решение: Буферизировать или правильно синхронизировать
//
// 5. Паника в processTask роняет воркер (строка 83)
//    Проблема: if task.Data == "corrupt" { panic("corrupted data!") }
//              Паника убивает воркер, он перестает обрабатывать задачи
//              Нет recover, паника убьет горутину
//    Решение: Добавить recover в воркер, обработать панику как ошибку
//
// 6. Нет таймаута на обработку задачи
//    Проблема: Задача может зависнуть навсегда
//    Решение: Использовать context.WithTimeout
//
// 7. Нет graceful shutdown
//    Проблема: Программа завершается сразу после обработки
//              Воркеры могут быть убиты в середине работы
//    Решение: Использовать WaitGroup и корректное завершение
//
// 8. Нет механизма отмены
//    Проблема: Нельзя остановить обработку извне
//    Решение: Использовать context для отмены
//
// 9. Результаты обрабатываются синхронно (строка 54-61)
//    Проблема: Главная горутина ждет каждый результат по очереди
//              Если воркеров меньше задач, может быть медленно
//    Решение: Использовать отдельную горутину для сбора результатов
//
// 10. Нет логирования и метрик
//     Проблема: Не видно что происходит с воркерами
//     Решение: Добавить логирование старта/стопа воркеров

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

type Task struct {
	ID   int
	Data string
}

type Result struct {
	TaskID int
	Output string
	Error  error
}

// WorkerPool управляет пулом воркеров
type WorkerPool struct {
	numWorkers int
	tasks      chan Task
	results    chan Result
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewWorkerPool создает новый пул воркеров
func NewWorkerPool(numWorkers int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	return &WorkerPool{
		numWorkers: numWorkers,
		tasks:      make(chan Task, numWorkers*2), // Буфер для smooth работы
		results:    make(chan Result, numWorkers*2),
		ctx:        ctx,
		cancel:     cancel,
	}
}

// Start запускает воркеры
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// worker обрабатывает задачи из канала
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	fmt.Printf("Worker %d started\n", id)
	defer fmt.Printf("Worker %d stopped\n", id)

	for {
		select {
		case task, ok := <-wp.tasks:
			if !ok {
				// Канал закрыт, завершаем воркер
				return
			}

			// Обрабатываем задачу с recover для перехвата паник
			wp.processTaskSafely(task)

		case <-wp.ctx.Done():
			// Контекст отменен, завершаем воркер
			return
		}
	}
}

// processTaskSafely обрабатывает задачу с защитой от паник
func (wp *WorkerPool) processTaskSafely(task Task) {
	// Используем defer + recover для перехвата паник
	defer func() {
		if r := recover(); r != nil {
			// Паника произошла, отправляем как ошибку
			wp.results <- Result{
				TaskID: task.ID,
				Error:  fmt.Errorf("panic: %v", r),
			}
		}
	}()

	// Создаем контекст с таймаутом для задачи
	ctx, cancel := context.WithTimeout(wp.ctx, 5*time.Second)
	defer cancel()

	// Обрабатываем задачу
	output, err := processTaskWithContext(ctx, task)

	// Отправляем результат
	select {
	case wp.results <- Result{
		TaskID: task.ID,
		Output: output,
		Error:  err,
	}:
	case <-wp.ctx.Done():
		// Контекст отменен, не отправляем результат
	}
}

// Submit отправляет задачу в пул
func (wp *WorkerPool) Submit(task Task) error {
	select {
	case wp.tasks <- task:
		return nil
	case <-wp.ctx.Done():
		return fmt.Errorf("worker pool is shutting down")
	}
}

// Results возвращает канал результатов
func (wp *WorkerPool) Results() <-chan Result {
	return wp.results
}

// Shutdown корректно завершает пул воркеров
func (wp *WorkerPool) Shutdown() {
	// Закрываем канал задач - воркеры завершат после обработки текущих
	close(wp.tasks)

	// Ждем завершения всех воркеров
	wp.wg.Wait()

	// Закрываем канал результатов
	close(wp.results)
}

// Stop немедленно останавливает все воркеры
func (wp *WorkerPool) Stop() {
	// Отменяем контекст - воркеры остановятся сразу
	wp.cancel()

	// Ждем завершения
	wp.wg.Wait()

	// Закрываем каналы
	close(wp.results)
}

// processTaskWithContext обрабатывает задачу с поддержкой отмены
func processTaskWithContext(ctx context.Context, task Task) (string, error) {
	// Канал для результата
	resultChan := make(chan string, 1)
	errChan := make(chan error, 1)

	// Запускаем обработку в отдельной горутине
	go func() {
		// Симулируем обработку
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

		// Некоторые задачи вызывают панику
		if task.Data == "corrupt" {
			panic("corrupted data!")
		}

		resultChan <- fmt.Sprintf("processed %s", task.Data)
	}()

	// Ждем результат или таймаут
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errChan:
		return "", err
	case <-ctx.Done():
		return "", fmt.Errorf("task timeout: %w", ctx.Err())
	}
}

// РЕШЕНИЕ 1: Простое использование пула
func solution1() {
	fmt.Println("=== Решение 1: Базовый пул воркеров ===")

	tasks := []Task{
		{ID: 1, Data: "image1.jpg"},
		{ID: 2, Data: "image2.jpg"},
		{ID: 3, Data: "corrupt"},
		{ID: 4, Data: "image4.jpg"},
		{ID: 5, Data: "image5.jpg"},
	}

	// Создаем пул с количеством воркеров = CPU cores
	numWorkers := runtime.NumCPU()
	pool := NewWorkerPool(numWorkers)

	// Запускаем воркеры
	pool.Start()

	// Отправляем задачи
	go func() {
		for _, task := range tasks {
			if err := pool.Submit(task); err != nil {
				fmt.Printf("Failed to submit task %d: %v\n", task.ID, err)
			}
		}
		// Завершаем прием задач
		pool.Shutdown()
	}()

	// Собираем результаты
	for result := range pool.Results() {
		if result.Error != nil {
			fmt.Printf("❌ Task %d failed: %v\n", result.TaskID, result.Error)
		} else {
			fmt.Printf("✅ Task %d completed: %s\n", result.TaskID, result.Output)
		}
	}

	fmt.Println()
}

// РЕШЕНИЕ 2: С таймаутом на весь пул
func solution2() {
	fmt.Println("=== Решение 2: С общим таймаутом ===")

	tasks := []Task{
		{ID: 1, Data: "image1.jpg"},
		{ID: 2, Data: "image2.jpg"},
		{ID: 3, Data: "image3.jpg"},
	}

	pool := NewWorkerPool(2)
	pool.Start()

	// Устанавливаем таймаут на всю обработку
	timeout := time.After(3 * time.Second)

	// Отправляем задачи
	go func() {
		for _, task := range tasks {
			pool.Submit(task)
		}
		pool.Shutdown()
	}()

	// Собираем результаты с таймаутом
	processed := 0
	for {
		select {
		case result, ok := <-pool.Results():
			if !ok {
				fmt.Printf("All tasks processed: %d\n", processed)
				return
			}
			processed++
			if result.Error != nil {
				fmt.Printf("❌ Task %d failed: %v\n", result.TaskID, result.Error)
			} else {
				fmt.Printf("✅ Task %d completed: %s\n", result.TaskID, result.Output)
			}
		case <-timeout:
			fmt.Println("⏱️ Timeout! Stopping pool...")
			pool.Stop()
			fmt.Printf("Processed %d tasks before timeout\n", processed)
			return
		}
	}
}

func main() {
	solution1()
	solution2()
}

// ПОДРОБНОЕ ОБЪЯСНЕНИЕ ПРОБЛЕМ:
//
// ПРОБЛЕМА 1: Неограниченное создание горутин
//
// Оригинальный код:
//   for range tasks {
//       go worker(jobChan, resultChan)
//   }
//
// Если tasks = 10000, создастся 10000 горутин!
// Проблемы:
// - Overhead на создание/управление горутинами
// - Переключение контекста между горутинами
// - Память (каждая горутина ~2KB стека)
// - Scheduler pressure
//
// Решение:
//   numWorkers := runtime.NumCPU()
//   for i := 0; i < numWorkers; i++ {
//       go worker()
//   }
//
// ПРОБЛЕМА 2: Goroutine leak
//
// Оригинальный код:
//   go func() {
//       for _, task := range tasks {
//           jobChan <- task
//       }
//   }()  // Канал НЕ закрывается!
//
//   func worker(jobs <-chan Task, ...) {
//       for job := range jobs {  // Ждет вечно!
//           ...
//       }
//   }
//
// Что происходит:
// 1. Все задачи отправлены
// 2. Главная горутина завершается
// 3. Воркеры ждут в range jobs (канал не закрыт)
// 4. Воркеры остаются в памяти навсегда (leak!)
//
// Как обнаружить:
//   runtime.NumGoroutine() - количество живых горутин
//   До: 1 (main)
//   После запуска: 6 (main + 5 workers)
//   После завершения: 6 (воркеры не завершились!)
//
// Решение:
//   close(jobChan)  // После отправки всех задач
//   // Теперь range jobs завершится
//
// ПРОБЛЕМА 3: Паника убивает воркер
//
// Оригинальный код:
//   func processTask(task Task) (string, error) {
//       if task.Data == "corrupt" {
//           panic("corrupted data!")  // Убивает горутину!
//       }
//   }
//
// Что происходит:
// T=0ms:   5 воркеров активны
// T=50ms:  Worker 3 обрабатывает "corrupt"
// T=50ms:  panic! Worker 3 умирает
// T=51ms:  Только 4 воркера работают
// T=100ms: Остальные задачи обрабатываются медленнее
//
// Решение:
//   defer func() {
//       if r := recover(); r != nil {
//           results <- Result{Error: fmt.Errorf("panic: %v", r)}
//       }
//   }()
//
// ПРОБЛЕМА 4: Нет graceful shutdown
//
// Оригинальный код завершается так:
// 1. Все результаты получены
// 2. main() завершается
// 3. Программа убивает все горутины (forceful kill)
// 4. Воркеры могут быть в середине обработки
//
// Правильный shutdown:
// 1. Перестаем принимать новые задачи
// 2. Закрываем канал задач
// 3. Воркеры завершают текущие задачи
// 4. Воркеры выходят из цикла (range завершается)
// 5. wg.Wait() дожидается всех воркеров
// 6. Закрываем канал результатов
// 7. Главная горутина завершается
//
// TEMPORAL DIAGRAM:
//
// Time    Main            Worker1         Worker2         Worker3
// ────────────────────────────────────────────────────────────────
// 0ms     Start pool
// 1ms     Submit Task1    ← Task1
// 2ms     Submit Task2                    ← Task2
// 3ms     Submit Task3                                    ← Task3
// 4ms     close(tasks)
// 5ms     wait...         Processing...   Processing...   Processing...
// 50ms    ← Result1       Done → nil
// 51ms                    Exit ✓
// 60ms    ← Result2                       Done → nil
// 61ms                                    Exit ✓
// 70ms    ← Result3                                       Done → nil
// 71ms                                                    Exit ✓
// 72ms    wg.Wait() ✓
// 73ms    close(results)
// 74ms    Exit ✓
//
// БЕЗ graceful shutdown:
// 70ms    Exit main()     KILLED!         KILLED!         KILLED!
//         ↑ Задачи теряются, воркеры убиты в середине работы
//
// БУФЕРИЗАЦИЯ КАНАЛОВ:
//
// Без буфера:
//   jobChan := make(chan Task)
//   // Каждая отправка блокируется до чтения
//   jobChan <- task1  // Блокируется пока воркер не прочитает
//
// С буфером:
//   jobChan := make(chan Task, 10)
//   // Можно отправить 10 задач без блокировки
//   jobChan <- task1  // OK
//   jobChan <- task2  // OK
//   ...
//   jobChan <- task10 // OK
//   jobChan <- task11 // Блокируется
//
// Размер буфера:
// - Слишком маленький: отправитель блокируется часто
// - Слишком большой: лишняя память
// - Оптимально: numWorkers * 2 (небольшая очередь)
//
// ЛУЧШИЕ ПРАКТИКИ:
//
// 1. Фиксированное количество воркеров (runtime.NumCPU())
// 2. Закрывать канал задач после отправки всех
// 3. Использовать recover() в воркерах
// 4. WaitGroup для синхронизации завершения
// 5. Context для отмены и таймаутов
// 6. Graceful shutdown: close() → wg.Wait() → close()
// 7. Буферизировать каналы для производительности
// 8. Логировать старт/стоп воркеров
// 9. Метрики: количество активных воркеров, обработанных задач
// 10. Тесты на goroutine leaks
