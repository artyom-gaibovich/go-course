package main

import (
	"fmt"
	"sync"
)

var (
	value    int
	sharedMu sync.Mutex
)

// BAD: a local mutex gives every goroutine its own copy -> no mutual exclusion.
//
//	func increment() { var mu sync.Mutex; mu.Lock(); value++; mu.Unlock() }
//
// GOOD: the mutex is shared between all goroutines.
func increment() {
	sharedMu.Lock()
	value++
	sharedMu.Unlock()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()
			increment()
		}()
	}
	wg.Wait()
	fmt.Println("value:", value)
}
