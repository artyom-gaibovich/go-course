/*
Домашнее задание 3: Охота на гонки.

Ниже три заведомо сломанных фрагмента. Почините КАЖДЫЙ и подтвердите
результат детектором гонок: go run -race main.go

Bug 1 (brokenLocalMutex): мьютекс объявлен ЛОКАЛЬНО внутри горутины —
       у каждой горутины свой экземпляр, синхронизации нет.

Bug 2 (LeakyBuffer.Data): метод отдаёт внутренний слайс наружу, минуя
       мьютекс — пользователь создаёт data race. Сделайте deep copy
       или функциональный обходчик ForEach.

Bug 3 (deferInLoop): defer mu.Unlock() в теле цикла — defer срабатывает
       на выходе из ФУНКЦИИ, а не итерации; второй Lock зависнет.
*/

package main

import (
	"fmt"
	"sync"
)

// --- Bug 1 ---
var counter int

func brokenLocalMutex() {
	var wg sync.WaitGroup
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()
			var mu sync.Mutex // TODO: вынести наружу как общий мьютекс
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}
	wg.Wait()
}

// --- Bug 2 ---
type LeakyBuffer struct {
	mu   sync.Mutex
	data []int
}

func (b *LeakyBuffer) Append(v int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.data = append(b.data, v)
}

func (b *LeakyBuffer) Data() []int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.data // TODO: не отдавать внутренний слайс — deep copy
}

// --- Bug 3 ---
var loopMu sync.Mutex

func deferInLoop(items []int) {
	for _, v := range items {
		loopMu.Lock()
		defer loopMu.Unlock() // TODO: defer сработает в конце функции — починить
		_ = v
	}
}

func main() {
	brokenLocalMutex()
	fmt.Println("counter:", counter)

	b := &LeakyBuffer{}
	b.Append(1)
	b.Append(2)
	fmt.Println("data:", b.Data())

	deferInLoop([]int{1, 2, 3})
	fmt.Println("done")
}
