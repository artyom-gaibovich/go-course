package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// generate пишет числа 0..n−1 в возвращаемый канал.
// Гарантия: канал закрывается когда все числа отправлены ИЛИ ctx отменён.
// Гарантия: после возврата из generate все горутины внутри завершены.
//
// ЗАДАЧА: реализуй generate.
// Требования:
// – Возвращает <-chan int (читатель не может писать в него)
// – Внутри горутина которая пишет числа и закрывает канал
// – Если ctx отменён до отправки — горутина выходит и закрывает канал
// – Буфер канала: 0 (небуферизованный)
func generate(ctx context.Context, n int) <-chan int {
	ch := make(chan int)

	go func(ctx context.Context) {
		defer close(ch)
		for i := range n {
			ch <- i
		}
	}(ctx)
	// TODO
	return ch
}

// process читает из in, обрабатывает каждый элемент параллельно (numWorkers горутин),
// пишет результаты в возвращаемый канал.
//
// Гарантия: выходной канал закрывается когда ВСЕ воркеры завершились.
// Гарантия: воркер завершается если: in закрыт ИЛИ ctx отменён.
// Гарантия: каждый элемент взятый из in будет обработан до конца —
//
//	даже если ctx отменился в процессе (но с учётом individualTimeout).
//
// ЗАДАЧА: реализуй process.
// Требования:
// – numWorkers горутин читают из in через общий канал (не делить диапазон)
// – Каждый элемент обрабатывается через work(ctx, v) с individualTimeout = 300ms
//
//	то есть: каждый вызов work получает context с таймаутом 300ms,
//	независимо от состояния внешнего ctx
//
// – Результат пишется в out — но если внешний ctx отменён и out заполнен,
//
//	воркер НЕ блокируется, а дропает результат
//
// – Выходной канал буферизованный: numWorkers
func process(ctx context.Context, in <-chan int, numWorkers int) <-chan int {
	out := make(chan int, numWorkers)

	wg := sync.WaitGroup{}
	wg.Add(numWorkers)

	for _ = range numWorkers {
		go func() {
			defer wg.Done()
			for {
				select {
				case i, ok := <-in:
					if !ok {
						return
					}
					res := work(ctx, i)
					select {
					case out <- res:
					case <-ctx.Done():
						return
					}
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// aggregate читает из in и собирает батчи размером batchSize.
// Когда батч заполнен ИЛИ прошло flushInterval — флашит батч вызовом flush(batch).
// Когда in закрыт — флашит остаток (даже если батч пустой — не флашить).
//
// Гарантия: aggregate блокирует вызывающего до полного завершения.
// Гарантия: flush вызывается последовательно (не параллельно).
//
// ЗАДАЧА: реализуй aggregate.
// Требования:
// – Один for-select по: in, ticker, ctx.Done()
// – При ctx.Done(): дочитать in до закрытия (drain), собрать в батч, сделать финальный flush
// – ticker должен быть остановлен (Stop) и его канал задренирован после использования
func aggregate(ctx context.Context, in <-chan int, batchSize int, flushInterval time.Duration, flush func([]int)) {
	batch := make([]int, 0, batchSize)
	ticker := time.NewTicker(flushInterval)
	defer ticker.Stop()

	for {
		select {
		case res, ok := <-in:
			if !ok {
				if len(batch) > 0 {
					flush(batch)
				}
				return
			}
			batch = append(batch, res)
			if len(batch) == batchSize {
				flush(batch)
				batch = make([]int, 0, batchSize)
			}
		case <-ctx.Done():
			for res := range in {
				batch = append(batch, res)
			}
			flush(batch)
		case <-ticker.C:
			if len(batch) == batchSize {
				flush(batch)
				batch = make([]int, 0, batchSize)
			}
		}
	}

}

// work имитирует тяжёлую обработку. Не меняй.
func work(ctx context.Context, v int) int {
	select {
	case <-ctx.Done():
		return -v
	case <-time.After(time.Duration(rand.Intn(400)) * time.Millisecond):
		return v * 2
	}
}

// Не меняй main.
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	src := generate(ctx, 50)
	processed := process(ctx, src, 6)

	var mu sync.Mutex
	var batches [][]int

	aggregate(ctx, processed, 5, 300*time.Millisecond, func(batch []int) {
		mu.Lock()
		batches = append(batches, append([]int{}, batch...))
		mu.Unlock()
	})

	total := 0
	for i, b := range batches {
		fmt.Printf("batch %d: %v\n", i, b)
		total += len(b)
	}
	fmt.Printf("total processed: %d\n", total)
}
