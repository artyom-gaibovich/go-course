package main

import (
	"fmt"
	"sync"
)

// Run with: go run -race main.go
// The race detector flags the DATA RACE on `msg`, but cannot see the
// RACE CONDITION: the final value depends on goroutine ordering.
//
// To fix the data race add a mutex; to fix the race condition impose an
// explicit order (here: Wait between the two goroutines).
func main() {
	var msg string
	var mu sync.Mutex
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		mu.Lock()
		msg = "hello world"
		mu.Unlock()
	}()
	wg.Wait() // order enforced: writer always runs before reader

	wg.Add(1)
	go func() {
		defer wg.Done()
		mu.Lock()
		fmt.Println("read:", msg)
		mu.Unlock()
	}()
	wg.Wait()
}
