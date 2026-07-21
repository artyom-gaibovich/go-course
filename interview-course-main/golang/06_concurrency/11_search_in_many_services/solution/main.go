package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// SearchFunc - функция поиска на одном сервере.
type SearchFunc func(server, query string) ([]string, error)

// Search выполняет параллельный поиск на нескольких серверах.
// Возвращает первый успешный результат или ошибку, если все серверы упали.
func Search(servers []string, query string, searchFunc SearchFunc) ([]string, error) {
	// Создаем буферизированные каналы для результатов и ошибок.
	// Буфер = len(servers) чтобы горутины не блокировались при отправке.
	resultChan := make(chan []string, len(servers))
	errorChan := make(chan error, len(servers))

	// WaitGroup для отслеживания завершения горутин
	var wg sync.WaitGroup
	wg.Add(len(servers))

	// Запускаем горутину для каждого сервера.
	for _, server := range servers {
		// Захватываем переменную для горутины
		// Актуально для версий Go до 1.22
		server := server
		go func() {
			defer wg.Done() // Уменьшаем счетчик при выходе из горутины

			// Выполняем поиск на конкретном сервере.
			result, err := searchFunc(server, query)
			if err != nil {
				// Отправляем ошибку в канал ошибок.
				errorChan <- err
				return
			}
			// Отправляем успешный результат.
			resultChan <- result
		}()
	}

	// Закрываем каналы после завершения всех горутин
	go func() {
		wg.Wait()
		close(resultChan)
		close(errorChan)
	}()

	// Собираем результаты от всех серверов.
	errCount := 0
	for i := 0; i < len(servers); i++ {
		select {
		case result, ok := <-resultChan:
			if !ok {
				// Канал закрыт, все горутины завершились
				break
			}
			// Получили первый успешный результат - возвращаем его.
			return result, nil
		case _, ok := <-errorChan:
			if !ok {
				// Канал закрыт
				break
			}
			// Получили ошибку, увеличиваем счетчик.
			errCount++
			// Если все серверы вернули ошибку.
			if errCount == len(servers) {
				return nil, errors.New("all servers failed")
			}
		}
	}

	return nil, errors.New("unexpected error")
}

// Пример использования
func main() {
	servers := []string{"server1", "server2", "server3"}

	// Симуляция поиска с задержками и ошибками.
	mockSearch := func(server, query string) ([]string, error) {
		delay := time.Duration(rand.Intn(300)) * time.Millisecond
		time.Sleep(delay)

		// server1 быстро возвращает ошибку
		if server == "server1" {
			return nil, fmt.Errorf("%s failed", server)
		}

		// Остальные возвращают результаты
		return []string{fmt.Sprintf("result from %s for %s", server, query)}, nil
	}

	result, err := Search(servers, "test query", mockSearch)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Success: %v\n", result)
}

// РЕШЕНИЕ:
//
// 1. Используем sync.WaitGroup:
//    var wg sync.WaitGroup
//    wg.Add(len(servers))
//
//    go func() {
//        defer wg.Done()
//        // работа горутины
//    }()
//
// 2. Закрываем каналы после завершения всех горутин:
//    go func() {
//        wg.Wait()           // Ждем все горутины
//        close(resultChan)   // Закрываем канал результатов
//        close(errorChan)    // Закрываем канал ошибок
//    }()
//
// 3. Проверяем закрытие каналов при чтении:
//    result, ok := <-resultChan
//    if !ok {
//        // Канал закрыт
//    }
//
// МЕХАНИЗМ РАБОТЫ:
//
// Временная шкала:
//   00:00.000 - Запускаем 3 горутины
//   00:00.000 - wg.Add(3), счетчик = 3
//   00:00.000 - Запускаем фоновую горутину для закрытия каналов
//   00:00.050 - server1 завершается, wg.Done(), счетчик = 2
//   00:00.100 - server2 завершается, wg.Done(), счетчик = 1
//   00:00.150 - server3 завершается, wg.Done(), счетчик = 0
//   00:00.150 - wg.Wait() в фоновой горутине разблокируется
//   00:00.150 - close(resultChan), close(errorChan)
//   00:00.150 - select в main видит закрытые каналы
//
// ПОЧЕМУ ЭТО ВАЖНО:
//
// 1. Корректное завершение:
//    Закрытие каналов сигнализирует читателям, что данных больше не будет.
//
// 2. Range по каналам:
//    for result := range resultChan {
//        // обработка
//    }
//    // Цикл завершится только когда канал закроется
//
// 3. Идиоматичность:
//    В Go принято закрывать каналы, когда отправитель завершил работу.
//    "Don't communicate by sharing memory; share memory by communicating."
//
// АЛЬТЕРНАТИВНЫЕ ПОДХОДЫ:
//
// 1. Context для отмены (если нужна отмена извне):
//    ctx, cancel := context.WithCancel(ctx)
//    defer cancel()
//    // Но для простого "first wins" это избыточно
//
// 2. Один канал вместо двух:
//    type Result struct {
//        Data []string
//        Err  error
//    }
//    resultChan := make(chan Result)
//    // Проще управление, но теряется type safety
//
// 3. Errgroup:
//    g, ctx := errgroup.WithContext(context.Background())
//    for _, server := range servers {
//        g.Go(func() error { ... })
//    }
//    err := g.Wait()
//    // Автоматически управляет горутинами и ошибками
//
// BEST PRACTICES:
//
// 1. Отправитель закрывает канал, не получатель:
//    ✅ Правильно: горутины отправляют, отдельная горутина закрывает
//    ❌ Неправильно: получатель закрывает канал
//
// 2. Закрывать канал только один раз:
//    Повторное закрытие канала вызывает panic.
//    Используйте sync.Once если есть риск двойного закрытия.
//
// 3. Буферизация для избежания блокировок:
//    resultChan := make(chan T, len(workers))
//    Горутины не блокируются при отправке, даже если никто не читает.
//
// 4. Проверка закрытия при чтении:
//    value, ok := <-ch
//    if !ok {
//        // Канал закрыт и пуст
//    }
//
// 5. WaitGroup для синхронизации:
//    Используйте WaitGroup для ожидания завершения группы горутин.
//
// КЛЮЧЕВЫЕ ВЫВОДЫ:
//
// - Всегда закрывайте каналы после завершения отправителей
// - Используйте WaitGroup для отслеживания горутин
// - Проверяйте ok при чтении из канала
// - Буферизация помогает избежать блокировок
// - Простое решение лучше сложного, если оно решает задачу
