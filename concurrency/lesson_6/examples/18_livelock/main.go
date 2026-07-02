package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Livelock: goroutines stay busy (active waiting) but make no progress.
// Each grabs one lock, yields, and keeps retrying the other, backing off
// when it can't get both. The runtime does NOT report a deadlock here —
// the goroutines are running, just doing useless work.
//
// We cap it with a timeout so the demo terminates instead of spinning forever.
func main() {
	var mu1, mu2 sync.Mutex
	var progress int64

	deadline := time.Now().Add(300 * time.Millisecond)

	worker := func(first, second *sync.Mutex) {
		for time.Now().Before(deadline) {
			first.Lock()
			if second.TryLock() {
				atomic.AddInt64(&progress, 1)
				second.Unlock()
				first.Unlock()
				return
			}
			first.Unlock()
			time.Sleep(time.Millisecond) // back off and retry: looks busy, no progress
		}
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); worker(&mu1, &mu2) }()
	go func() { defer wg.Done(); worker(&mu2, &mu1) }()
	wg.Wait()

	fmt.Println("progress made:", atomic.LoadInt64(&progress))
}
