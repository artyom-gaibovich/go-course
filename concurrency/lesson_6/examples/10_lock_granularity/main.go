package main

import (
	"fmt"
	"sync"
)

var (
	mu    sync.Mutex
	cache = map[string]string{"key": "value"}
)

// Narrow critical section: read into a local under the lock, release,
// then do the slow work (printing) without holding the mutex.
func doSomething(key string) {
	mu.Lock()
	item := cache[key]
	mu.Unlock()

	fmt.Println("processing:", item)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func() {
			defer wg.Done()
			doSomething("key")
		}()
	}
	wg.Wait()
}
