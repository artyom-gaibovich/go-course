package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type Result struct {
	Value string
}

// JobFunc - асинхронная задача.
type JobFunc func(ctx context.Context, input string) (Result, error)

// MultiProcess выполняет все jobs параллельно и возвращает первый успешный результат.
func MultiProcess(ctx context.Context, input string, jobs []JobFunc) (Result, error) {
	if len(jobs) == 0 {
		return Result{}, errors.New("no jobs provided")
	}

	// Каналы для результатов и ошибок.
	resultChan := make(chan Result, len(jobs))
	errorChan := make(chan error, len(jobs))

	// WaitGroup для отслеживания завершения горутин
	var wg sync.WaitGroup
	wg.Add(len(jobs))

	// Запускаем все jobs параллельно.
	for _, job := range jobs {
		// Захватываем переменную для горутины
		// Актуально для версий Go до 1.22
		job := job
		go func() {
			// Уменьшаем счетчик при выходе
			defer wg.Done()

			result, err := job(ctx, input)
			if err != nil {
				errorChan <- err
				return
			}
			resultChan <- result
		}()
	}

	// Закрываем каналы после завершения всех горутин
	go func() {
		wg.Wait()
		close(resultChan)
		close(errorChan)
	}()

	// Собираем результаты.
	errCount := 0
	var lastErr error

	for i := 0; i < len(jobs); i++ {
		select {
		case result, ok := <-resultChan:
			if !ok {
				// Канал закрыт
				break
			}
			// Первый успешный результат - возвращаем его.
			return result, nil
		case err, ok := <-errorChan:
			if !ok {
				// Канал закрыт
				break
			}
			// Собираем ошибки.
			errCount++
			lastErr = err
			if errCount == len(jobs) {
				// Все jobs вернули ошибки.
				return Result{}, fmt.Errorf("all jobs failed, last error: %w", lastErr)
			}
		case <-ctx.Done():
			// Контекст отменен.
			return Result{}, ctx.Err()
		}
	}

	return Result{}, errors.New("unexpected error")
}

func main() {
	// Пример использования.
	job1 := func(ctx context.Context, input string) (Result, error) {
		time.Sleep(200 * time.Millisecond)
		return Result{}, errors.New("job1 failed")
	}

	job2 := func(ctx context.Context, input string) (Result, error) {
		time.Sleep(100 * time.Millisecond)
		return Result{Value: "job2 success: " + input}, nil
	}

	job3 := func(ctx context.Context, input string) (Result, error) {
		time.Sleep(150 * time.Millisecond)
		return Result{Value: "job3 success: " + input}, nil
	}

	ctx := context.Background()
	result, err := MultiProcess(ctx, "test input", []JobFunc{job1, job2, job3})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Success: %+v\n", result)
}

// РЕШЕНИЕ:
//
// 1. Добавляем WaitGroup для отслеживания горутин:
//    var wg sync.WaitGroup
//    wg.Add(len(jobs))
//
//    go func() {
//        defer wg.Done()
//        // выполнение job
//    }()
//
// 2. Закрываем каналы после завершения всех jobs:
//    go func() {
//        wg.Wait()           // Ждем все горутины
//        close(resultChan)   // Закрываем каналы
//        close(errorChan)
//    }()
//
// 3. Проверяем закрытие при чтении:
//    result, ok := <-resultChan
//    if !ok {
//        // Канал закрыт и пуст
//    }
//
// МЕХАНИЗМ РАБОТЫ:
//
// Временная шкала:
//   00:00.000 - Запускаем 3 jobs
//   00:00.000 - wg.Add(3), счетчик = 3
//   00:00.000 - Запускаем фоновую горутину для закрытия
//   00:00.100 - job2 завершается, wg.Done(), счетчик = 2
//   00:00.100 - MultiProcess() получает результат job2
//   00:00.100 - return result, nil
//   00:00.150 - job3 завершается, wg.Done(), счетчик = 1
//   00:00.200 - job1 завершается, wg.Done(), счетчик = 0
//   00:00.200 - wg.Wait() разблокируется
//   00:00.200 - close(resultChan), close(errorChan)
//
// ВАЖНЫЙ МОМЕНТ:
//
// После return в MultiProcess() горутины продолжат работать:
// - job1 и job3 завершат выполнение
// - Они отправят свои результаты в буферизированные каналы
// - Фоновая горутина дождется их и закроет каналы
// - Все корректно завершится, без утечек
//
// ПОЧЕМУ БУФЕРИЗАЦИЯ ВАЖНА:
//
// Без буфера:
//   resultChan := make(chan Result)  // Unbuffered
//
//   // job2 завершается первым
//   resultChan <- result  // OK, main читает
//
//   // main возвращает результат
//   return result, nil
//
//   // job3 завершается
//   resultChan <- result  // БЛОКИРУЕТСЯ! Никто не читает
//   // Горутина зависает навсегда
//
// С буфером:
//   resultChan := make(chan Result, len(jobs))
//
//   // job3 завершается после return
//   resultChan <- result  // OK, помещается в буфер
//   // Горутина завершается
//
// ЛУЧШИЕ ПРАКТИКИ:
//
// 1. Всегда закрывайте каналы:
//    ✅ После завершения всех отправителей
//    ❌ Не закрывайте в получателе
//
// 2. Используйте WaitGroup:
//    Для отслеживания завершения группы горутин
//
// 3. Буферизация = количество отправителей:
//    resultChan := make(chan T, len(workers))
//    Предотвращает блокировки
//
// 4. Проверяйте закрытие:
//    value, ok := <-ch
//    if !ok { /* канал закрыт */ }
//
// 5. Один отправитель закрывает:
//    Не закрывайте канал из нескольких горутин!
//    Используйте отдельную горутину с wg.Wait()
//
// АЛЬТЕРНАТИВНЫЕ ПОДХОДЫ:
//
// 1. Context для отмены (если нужно):
//    В данной задаче ctx уже передается для других целей.
//    Можно использовать ctx.Done() для досрочной остановки jobs.
//
// 2. Errgroup:
//    import "golang.org/x/sync/errgroup"
//
//    g, ctx := errgroup.WithContext(ctx)
//    for _, job := range jobs {
//        g.Go(func() error { return job(ctx, input) })
//    }
//    err := g.Wait()
//
//    Автоматически управляет горутинами, но ждет все jobs.
//
// 3. Один канал для результатов и ошибок:
//    type JobResult struct {
//        Result Result
//        Error  error
//    }
//    ch := make(chan JobResult, len(jobs))
//
//    Проще, но теряется type safety.
//
// 4. Без каналов (callback):
//    var once sync.Once
//    var result Result
//
//    for _, job := range jobs {
//        go func() {
//            r, err := job(ctx, input)
//            if err == nil {
//                once.Do(func() { result = r })
//            }
//        }()
//    }
//
//    Но теперь нужно ждать wg.Wait() перед чтением result.
//
// СРАВНЕНИЕ ПОДХОДОВ:
//
// WaitGroup + close (используемый):
//   + Простая реализация
//   + Идиоматично для Go
//   + Все jobs завершаются
//   + Корректное управление ресурсами
//   - Нельзя отменить jobs досрочно
//
// Context cancellation:
//   + Можно отменить jobs после первого результата
//   + Экономия ресурсов
//   - Требует поддержки отмены в jobs
//   - Сложнее реализация
//
// Errgroup:
//   + Автоматическое управление
//   + Стандартная библиотека (x/sync)
//   - Ждет все jobs (не "first wins")
//   - Дополнительная зависимость
//
// КОГДА ИСПОЛЬЗОВАТЬ:
//
// WaitGroup + close:
//   - Простые учебные примеры
//   - Jobs быстро выполняются
//   - Нет требований к отмене
//
// Context:
//   - Долгие операции (HTTP, DB)
//   - Production код
//   - Нужна отмена для экономии ресурсов
//
// КЛЮЧЕВЫЕ ВЫВОДЫ:
//
// - Закрывайте каналы после завершения отправителей
// - WaitGroup помогает отследить завершение всех горутин
// - Буферизация = len(workers) предотвращает блокировки
// - Один отправитель (фоновая горутина) закрывает каналы
// - Простота лучше сложности, если решает задачу
