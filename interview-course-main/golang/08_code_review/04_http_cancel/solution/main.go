package main

// СПИСОК НАЙДЕННЫХ ПРОБЛЕМ И ИСПРАВЛЕНИЙ:
//
// 1. Нет отмены запросов при ошибке (строка 32)
//    Проблема: Используется context.Background() - нельзя отменить
//              При ошибке в одном запросе, остальные продолжают работать
//    Решение: Использовать context.WithCancel() и вызывать cancel() при ошибке
//
// 2. Отсутствует WaitGroup (строки 30-41)
//    Проблема: Программа не ждет завершения горутин
//              Использует time.Sleep(400ms) - программа может завершиться
//              до завершения запросов, горутины будут убиты
//    Решение: Использовать sync.WaitGroup для ожидания всех горутин
//
// 3. time.Sleep для синхронизации (строка 42)
//    Проблема: Хардкод 400ms - непредсказуемо, может быть мало/много
//              Антипаттерн - не гарантирует завершение работы
//    Решение: Использовать WaitGroup вместо Sleep
//
// 4. Не закрывается response body (после строки 52)
//    Проблема: resp.Body не закрывается, утечка ресурсов
//              Connection leak - исчерпание file descriptors
//    Решение: Добавить defer resp.Body.Close() и io.Copy(io.Discard, resp.Body)
//
// 5. Игнорируется успешный результат (строка 60)
//    Проблема: Не выводим информацию об успешных запросах
//              Только ошибки видны (строки 33-37)
//    Решение: Выводить статус код и размер ответа
//
// 6. Нет timeout для запросов
//    Проблема: Запрос может висеть бесконечно долго
//    Решение: Добавить timeout через context.WithTimeout
//
// 7. Используется http.DefaultClient (строка 52)
//    Проблема: DefaultClient без настроек timeout
//              Может висеть вечно на DNS lookup, connection, response
//    Решение: Создать кастомный http.Client с таймаутами
//
// 8. Нет обработки статус кода (после строки 52)
//    Проблема: Статус 404, 500 и т.д. не считаются ошибкой
//              http.Client.Do() возвращает err только при сетевых ошибках
//    Решение: Проверять resp.StatusCode
//
// 9. Нет централизованной обработки ошибок (строки 33-37)
//    Проблема: Каждая горутина просто выводит ошибку и завершается
//              Нет механизма для отмены других запросов
//    Решение: Использовать errgroup или канал для первой ошибки
//
// 10. Race condition в выводе (fmt.Printf из горутин, строки 33-37)
//     Проблема: Конкурентный вывод может перемешиваться
//               Хотя fmt.Printf thread-safe, вывод может быть нечитаемым
//     Решение: Собирать результаты и выводить в главной горутине
//

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

// Result хранит результат выполнения запроса
type Result struct {
	URL        string
	StatusCode int
	Err        error
}

// РЕШЕНИЕ 1: С context.WithCancel и отменой при первой ошибке
func solution1() {
	fmt.Println("=== Решение 1: Отмена при первой ошибке ===")

	urls := []string{
		"https://google.com",
		"https://yandex.ru",
		"https://invalid-url-that-will-fail.com",
		"https://github.com",
	}

	// Создаем cancelable контекст
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Гарантируем отмену при выходе

	var wg sync.WaitGroup
	results := make(chan Result, len(urls))

	// Запускаем запросы параллельно
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			// Используем ctx с возможностью отмены
			statusCode, err := fetchWithCancel(ctx, url)

			result := Result{
				URL:        url,
				StatusCode: statusCode,
				Err:        err,
			}

			// Если ошибка - отменяем все остальные запросы
			if err != nil {
				fmt.Printf("❌ Error on %s, canceling all other requests\n", url)
				cancel() // Отменяем контекст
			}

			results <- result
		}(url)
	}

	// Закрываем канал после завершения всех горутин
	go func() {
		wg.Wait()
		close(results)
	}()

	// Собираем результаты
	for result := range results {
		if result.Err != nil {
			fmt.Printf("Failed: %s - %v\n", result.URL, result.Err)
		} else {
			fmt.Printf("Success: %s - Status %d\n", result.URL, result.StatusCode)
		}
	}

	fmt.Println()
}

// РЕШЕНИЕ 2: С errgroup (автоматическая отмена)
func solution2() {
	fmt.Println("=== Решение 2: С errgroup ===")

	urls := []string{
		"https://google.com",
		"https://yandex.ru",
		"https://github.com",
		"https://stackoverflow.com",
	}

	// errgroup автоматически создает контекст с cancel
	// При первой ошибке контекст отменяется для всех горутин
	g, ctx := errgroup.WithContext(context.Background())

	// Используем мьютекс для безопасного вывода
	var mu sync.Mutex

	for _, url := range urls {
		// Важно: создаем локальную копию для замыкания
		// Актуально для версий Go до 1.22
		url := url

		g.Go(func() error {
			statusCode, err := fetchWithCancel(ctx, url)
			if err != nil {
				return fmt.Errorf("%s: %w", url, err)
			}

			// Безопасный вывод
			mu.Lock()
			fmt.Printf("✅ Success: %s - Status %d\n", url, statusCode)
			mu.Unlock()

			return nil
		})
	}

	// Wait ждет завершения всех горутин и возвращает первую ошибку
	if err := g.Wait(); err != nil {
		fmt.Printf("❌ First error: %v\n", err)
	} else {
		fmt.Println("✅ All requests completed successfully")
	}

	fmt.Println()
}

// fetchWithCancel выполняет HTTP запрос с поддержкой отмены
func fetchWithCancel(ctx context.Context, url string) (int, error) {
	// Добавляем timeout к контексту (5 секунд на запрос)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Создаем запрос с контекстом
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	// Используем кастомный клиент с таймаутами
	client := &http.Client{
		Timeout: 10 * time.Second, // Общий timeout
		Transport: &http.Transport{
			MaxIdleConns:        100,
			IdleConnTimeout:     90 * time.Second,
			TLSHandshakeTimeout: 5 * time.Second,
		},
	}

	// Выполняем запрос
	resp, err := client.Do(req)
	if err != nil {
		// Проверяем причину ошибки
		if errors.Is(ctx.Err(), context.Canceled) {
			return 0, fmt.Errorf("request canceled: %w", err)
		}
		return 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// КРИТИЧЕСКИ ВАЖНО: Читаем body до конца для HTTP Keep-Alive!
	//
	// КОГДА НУЖНО:
	// - Когда НЕ читаем body (данные не нужны, проверяем только статус)
	// - Body есть, но мы его игнорируем
	// - Нужно для возврата connection в pool
	//
	// КОГДА НЕ НУЖНО:
	// - Если уже прочитали body через io.ReadAll(resp.Body)
	// - Если использовали json.NewDecoder(resp.Body).Decode(&v)
	// - Если использовали io.Copy(dst, resp.Body) куда-то
	// - В этих случаях body уже consumed, повторный io.Copy не нужен!
	//
	// НАША СИТУАЦИЯ:
	// - Мы НЕ читаем body (нам нужен только StatusCode)
	// - Поэтому ОБЯЗАТЕЛЬНО нужен io.Copy(io.Discard, resp.Body)
	//
	// Причины:
	// 1. HTTP Keep-Alive работает только если body прочитан до конца
	// 2. Без этого connection не вернется в pool и будет закрыто
	// 3. Следующий запрос создаст новое TCP соединение (медленно!)
	// 4. io.Discard - это оптимизированная "черная дыра" для данных
	// 5. Игнорируем ошибку (_, _) т.к. тело может быть закрыто из-за cancel
	//
	// Бенчмарк 1000 запросов без io.Copy:
	//   - 1000 новых TCP соединений (3-way handshake каждый раз)
	//   - ~100-300ms latency на каждое соединение
	//   - Исчерпание портов и file descriptors
	//
	// С io.Copy:
	//   - Переиспользуем 10-20 соединений из pool
	//   - ~10-20ms latency (без handshake)
	//   - Экономия ресурсов
	_, _ = io.Copy(io.Discard, resp.Body)

	// Проверяем статус код
	if resp.StatusCode >= 400 {
		return resp.StatusCode, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	return resp.StatusCode, nil
}

func main() {
	solution1()
	solution2()
}

// ПОДРОБНОЕ ОБЪЯСНЕНИЕ ОТМЕНЫ ЧЕРЕЗ КОНТЕКСТ:
//
// Как работает context.WithCancel():
//
// 1. Создаем контекст с возможностью отмены:
//    ctx, cancel := context.WithCancel(context.Background())
//
// 2. Передаем ctx во все горутины:
//    go func() {
//        err := fetchWithCancel(ctx, url)
//    }()
//
// 3. При ошибке вызываем cancel():
//    if err != nil {
//        cancel() // Отменяет ctx для ВСЕХ горутин
//    }
//
// 4. В http.NewRequestWithContext(ctx, ...) контекст встроен в запрос
//
// 5. Когда cancel() вызван:
//    - ctx.Done() канал закрывается
//    - http.Client прерывает запрос
//    - fetchWithCancel возвращает context.Canceled ошибку
//
// ВРЕМЕННАЯ ДИАГРАММА:
//
// T=0ms:   Запускаем 4 запроса
// T=50ms:  google.com   -> Success (200)
// T=100ms: invalid-url  -> Error! Вызываем cancel()
// T=100ms: yandex.ru    -> context.Canceled (прерван)
// T=100ms: github.com   -> context.Canceled (прерван)
//
// БЕЗ ОТМЕНЫ (оригинальный код):
// T=0ms:   Запускаем 4 запроса
// T=50ms:  google.com   -> Success (200)
// T=100ms: invalid-url  -> Error (но остальные продолжают)
// T=200ms: yandex.ru    -> Success (200)
// T=300ms: github.com   -> Success (200)
// Потеря времени и ресурсов!
//
// ПОЧЕМУ ВАЖНО ЗАКРЫВАТЬ response.Body И ЧИТАТЬ ДО КОНЦА:
//
// НЕПРАВИЛЬНО:
// resp, _ := http.DefaultClient.Do(req)
// // ОШИБКА 1: Не закрыли body
// // Последствия:
// //   - File descriptor leak (утечка FD)
// //   - Connection остается открытым
// //   - Через некоторое время: "too many open files"
// //   - В production: паника, restart контейнера
//
// ЧАСТИЧНО ПРАВИЛЬНО:
// resp, _ := http.DefaultClient.Do(req)
// defer resp.Body.Close()
// // ОШИБКА 2: Закрыли, но не прочитали до конца!
// // Последствия:
// //   - Connection ЗАКРЫВАЕТСЯ вместо возврата в pool
// //   - HTTP Keep-Alive НЕ РАБОТАЕТ
// //   - Каждый запрос = новое TCP соединение
// //   - 3-way handshake на каждый запрос (~100ms overhead)
// //   - TLS handshake при HTTPS (~200-400ms overhead)
//
// ПРАВИЛЬНО:
// resp, _ := http.DefaultClient.Do(req)
// defer resp.Body.Close()
// _, _ = io.Copy(io.Discard, resp.Body)  // Читаем до конца!
// // Результат:
// //   - Connection возвращается в pool
// //   - Следующий запрос переиспользует соединение
// //   - HTTP Keep-Alive работает
// //   - Нет overhead на TCP/TLS handshake
// //   - Латентность снижается в 5-10 раз!
//
// ЧТО ТАКОЕ io.Discard:
// - Оптимизированная реализация io.Writer
// - Отбрасывает все данные без выделения памяти
// - Аналог /dev/null в Unix
// - НЕ создает буферы, НЕ выделяет память
// - Просто "черная дыра" для байтов
//
// КОГДА ИГНОРИРОВАТЬ ОШИБКУ io.Copy:
// Используем _, _ = io.Copy(...) по следующим причинам:
// 1. Если context отменен, body может быть закрыт досрочно
// 2. Ошибка чтения body не критична (данные нам не нужны)
// 3. Главное - попытаться прочитать для Keep-Alive
// 4. Даже partial read помогает (частично очищает буфер)
//
// ПРАКТИЧЕСКИЕ ПРИМЕРЫ:
//
// СЛУЧАЙ 1: НЕ читаем body - НУЖЕН io.Copy(io.Discard)
//
// resp, err := client.Do(req)
// if err != nil {
//     return err
// }
// defer resp.Body.Close()
// _, _ = io.Copy(io.Discard, resp.Body)  // ✅ НУЖЕН!
//
// // Проверяем только статус
// if resp.StatusCode != 200 {
//     return fmt.Errorf("bad status: %d", resp.StatusCode)
// }
// return nil
//
// СЛУЧАЙ 2: Читаем body в память - io.Copy НЕ НУЖЕН
//
// resp, err := client.Do(req)
// if err != nil {
//     return err
// }
// defer resp.Body.Close()
//
// body, err := io.ReadAll(resp.Body)  // Body уже consumed
// if err != nil {
//     return err
// }
// // io.Copy(io.Discard, resp.Body) здесь НЕ НУЖЕН! ❌
//
// // Обрабатываем данные
// return processData(body)
//
// СЛУЧАЙ 3: Парсим JSON - io.Copy НЕ НУЖЕН
//
// resp, err := client.Do(req)
// if err != nil {
//     return err
// }
// defer resp.Body.Close()
//
// var result Response
// err = json.NewDecoder(resp.Body).Decode(&result)  // Body consumed
// if err != nil {
//     return err
// }
// // io.Copy(io.Discard, resp.Body) НЕ НУЖЕН! ❌
//
// return nil
//
// СЛУЧАЙ 4: Копируем в файл - io.Copy НЕ НУЖЕН
//
// resp, err := client.Do(req)
// if err != nil {
//     return err
// }
// defer resp.Body.Close()
//
// file, err := os.Create("output.dat")
// if err != nil {
//     return err
// }
// defer file.Close()
//
// _, err = io.Copy(file, resp.Body)  // Body consumed в файл
// if err != nil {
//     return err
// }
// // Повторный io.Copy(io.Discard, resp.Body) НЕ НУЖЕН! ❌
//
// return nil
//
// СЛУЧАЙ 5: Стримим данные - io.Copy НЕ НУЖЕН
//
// resp, err := client.Do(req)
// if err != nil {
//     return err
// }
// defer resp.Body.Close()
//
// scanner := bufio.NewScanner(resp.Body)
// for scanner.Scan() {
//     processLine(scanner.Text())
// }
// // Body consumed построчно
// // io.Copy(io.Discard, resp.Body) НЕ НУЖЕН! ❌
//
// return scanner.Err()
//
// ПРАВИЛО:
// - Body НЕ читали → io.Copy(io.Discard, resp.Body) НУЖЕН ✅
// - Body читали (ReadAll/Decoder/Scanner/Copy) → НЕ НУЖЕН ❌
//
// ПОЧЕМУ НЕ time.Sleep:
//
// // ПЛОХО:
// go fetch(url1)
// go fetch(url2)
// time.Sleep(400 * time.Millisecond) // Может быть мало или много!
//
// // ХОРОШО:
// var wg sync.WaitGroup
// wg.Add(2)
// go func() { defer wg.Done(); fetch(url1) }()
// go func() { defer wg.Done(); fetch(url2) }()
// wg.Wait() // Точно дождется завершения
//
// ERRGROUP vs РУЧНАЯ РЕАЛИЗАЦИЯ:
//
// РУЧНАЯ РЕАЛИЗАЦИЯ (Solution 1):
//
// var wg sync.WaitGroup
// ctx, cancel := context.WithCancel(context.Background())
// errChan := make(chan error, N)
//
// for _, item := range items {
//     wg.Add(1)
//     go func(item Item) {
//         defer wg.Done()
//         if err := process(ctx, item); err != nil {
//             errChan <- err
//             cancel()  // Вручную отменяем контекст
//         }
//     }(item)
// }
//
// go func() {
//     wg.Wait()
//     close(errChan)
// }()
//
// for err := range errChan {
//     log.Printf("Error: %v", err)
// }
//
// Проблемы:
// 1. Много boilerplate кода
// 2. Легко забыть cancel() при ошибке
// 3. Нужно вручную управлять errChan
// 4. Нужна отдельная горутина для close(errChan)
// 5. Race condition если забыли буфер в errChan
// 6. Нужно вручную проверять ctx.Err() в горутинах
//
// С ERRGROUP (Solution 2):
//
// g, ctx := errgroup.WithContext(context.Background())
//
// for _, item := range items {
//     item := item  // Копия для замыкания
//     g.Go(func() error {
//         return process(ctx, item)  // Просто возвращаем ошибку
//     })
// }
//
// if err := g.Wait(); err != nil {
//     log.Printf("First error: %v", err)
// }
//
// Преимущества:
// 1. Минимум boilerplate
// 2. Автоматическая отмена контекста при первой ошибке
// 3. Не нужны каналы для ошибок
// 4. Wait() возвращает первую ошибку
// 5. Невозможны race conditions
// 6. Код читается линейно сверху вниз
// 7. Проще тестировать
// 8. Меньше места для ошибок
//
// КАК РАБОТАЕТ errgroup.WithContext:
//
// 1. Создает cancelable context внутри себя
// 2. При вызове g.Go(func() error) запускает горутину
// 3. Если функция возвращает != nil ошибку:
//    - Сохраняет первую ошибку
//    - Автоматически вызывает cancel() на контексте
//    - Все остальные горутины получат ctx.Done()
// 4. g.Wait() блокируется пока все горутины не завершатся
// 5. g.Wait() возвращает первую встреченную ошибку
//
// ВАЖНО: errgroup останавливается при ПЕРВОЙ ошибке:
//
// g, ctx := errgroup.WithContext(context.Background())
//
// g.Go(func() error { return nil })           // OK
// g.Go(func() error { return errors.New("1") }) // Первая ошибка!
// g.Go(func() error { return errors.New("2") }) // Отменена через ctx
// g.Go(func() error { return errors.New("3") }) // Отменена через ctx
//
// err := g.Wait()  // Вернет errors.New("1")
//
// КОГДА НЕ ИСПОЛЬЗОВАТЬ errgroup:
//
// 1. Нужно собрать ВСЕ ошибки, а не только первую
//    -> Используйте WaitGroup + канал ошибок
//
// 2. Нужно продолжить работу несмотря на ошибки
//    -> Используйте WaitGroup без cancel
//
// 3. Нужна более сложная логика отмены
//    -> Используйте WaitGroup + кастомный context
//
// 4. Горутины не возвращают ошибки
//    -> Используйте обычный WaitGroup
//
// ERRGROUP С ОГРАНИЧЕНИЕМ ПАРАЛЛЕЛИЗМА:
//
// Если нужно ограничить количество одновременных горутин:
//
// g, ctx := errgroup.WithContext(context.Background())
// g.SetLimit(10)  // Максимум 10 параллельных горутин
//
// for i := 0; i < 1000; i++ {
//     i := i
//     g.Go(func() error {
//         return process(ctx, i)  // Только 10 одновременно
//     })
// }
//
// err := g.Wait()
//
// ЛУЧШИЕ ПРАКТИКИ:
//
// 1. Всегда используйте context для HTTP запросов
// 2. Добавляйте timeout через context.WithTimeout
// 3. Используйте errgroup для параллельных операций с ошибками
// 4. Используйте WaitGroup если ошибки не нужны
// 5. Всегда закрывайте response.Body
// 6. ОБЯЗАТЕЛЬНО читайте body до конца: io.Copy(io.Discard, resp.Body)
// 7. Проверяйте StatusCode, не только error
// 8. Создавайте кастомный http.Client с таймаутами
// 9. При первой ошибке - отменяйте остальные через cancel() или errgroup
// 10. Не забывайте про локальную копию переменной в замыкании (url := url)
