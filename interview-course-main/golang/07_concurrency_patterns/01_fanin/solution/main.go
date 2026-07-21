package main

import "sync"

// merge объединяет несколько каналов в один (fan-in паттерн).
func merge(channels ...chan int64) <-chan int64 {
	out := make(chan int64)
	var wg sync.WaitGroup

	// Запускаем горутину для каждого входного канала.
	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan int64) {
			defer wg.Done()
			// Читаем все значения из канала и отправляем в выходной.
			for v := range c {
				out <- v
			}
		}(ch)
	}

	// Закрываем выходной канал после завершения всех горутин.
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	channels := make([]chan int64, 10)
	for i := range channels {
		channels[i] = make(chan int64)
	}

	for i := range channels {
		go func(i int) {
			channels[i] <- int64(i)
			close(channels[i])
		}(i)
	}

	for v := range merge(channels...) {
		println(v)
	}
}

// Объяснение:
// 1. Fan-in паттерн объединяет несколько каналов в один.
// 2. Для каждого входного канала запускаем горутину, которая читает значения.
// 3. Все значения отправляются в один выходной канал.
// 4. Используем WaitGroup для ожидания завершения всех горутин.
// 5. Закрываем выходной канал после завершения всех чтений.
