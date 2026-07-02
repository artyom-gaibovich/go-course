package main

import (
	"fmt"
	"sync"
)

// withLock takes the mutex BY POINTER (never copy a sync primitive),
// runs action under Lock, and releases via defer so a panic inside
// action still unlocks the mutex.
func withLock(mu *sync.Mutex, action func()) {
	if action == nil {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	action()
}

var (
	mu      sync.Mutex
	counter int
)

func main() {
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			withLock(&mu, func() {
				counter++
			})
		}()
	}
	wg.Wait()
	fmt.Println("counter:", counter)
}
