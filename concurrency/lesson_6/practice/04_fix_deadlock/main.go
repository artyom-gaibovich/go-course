/*
Задание 4 (со звёздочкой): Найти и починить deadlock.

Код ниже иногда зависает с сообщением:
  fatal error: all goroutines are asleep - deadlock!

Запустите его несколько раз (или в цикле в шелле), поймайте deadlock,
нарисуйте на бумаге диаграмму Холта и почините.

Подсказка: проблема не в количестве мьютексов, а в ПОРЯДКЕ их захвата
в разных горутинах. Обеспечьте согласованный (одинаковый) порядок.

Сейчас normalizeAB и normalizeBA берут мьютексы в РАЗНОМ порядке — это баг.
*/

package main

import (
	"fmt"
	"sync"
)

var (
	muA sync.Mutex
	muB sync.Mutex
)

func normalizeAB() {
	muA.Lock()
	muB.Lock()
	muB.Unlock()
	muA.Unlock()
}

func normalizeBA() {
	// TODO: починить порядок — брать сначала muA, потом muB
	muB.Lock()
	muA.Lock()
	muA.Unlock()
	muB.Unlock()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2000)
	for i := 0; i < 1000; i++ {
		go func() { defer wg.Done(); normalizeAB() }()
		go func() { defer wg.Done(); normalizeBA() }()
	}
	wg.Wait()
	fmt.Println("done")
}
