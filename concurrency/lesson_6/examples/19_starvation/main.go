package main

import (
	"fmt"
	"sync"
	"time"
)

// The greedy worker holds the mutex for long stretches; the polite worker
// needs it briefly but often and ends up starved — far fewer iterations.
func main() {
	var mu sync.Mutex
	var wg sync.WaitGroup
	runtime := time.Second

	greedy := func() {
		defer wg.Done()
		var count int
		stop := time.Now().Add(runtime)
		for time.Now().Before(stop) {
			mu.Lock()
			time.Sleep(3 * time.Microsecond) // holds the lock a long time
			mu.Unlock()
			count++
		}
		fmt.Println("greedy iterations:", count)
	}

	polite := func() {
		defer wg.Done()
		var count int
		stop := time.Now().Add(runtime)
		for time.Now().Before(stop) {
			mu.Lock()
			time.Sleep(time.Microsecond) // holds it briefly
			mu.Unlock()
			count++
		}
		fmt.Println("polite iterations:", count)
	}

	wg.Add(2)
	go greedy()
	go polite()
	wg.Wait()
}
