package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var mu sync.Mutex
	var wg sync.WaitGroup

	mu.Lock()
	fmt.Println("main holds the mutex")

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("goroutine: trying to Lock -> goes into wait queue (waiting)")
		mu.Lock()
		fmt.Println("goroutine: got the mutex after main released it")
		mu.Unlock()
	}()

	time.Sleep(100 * time.Millisecond)
	fmt.Println("main releases the mutex, waiting goroutine becomes ready")
	mu.Unlock()

	wg.Wait()
}
