package main

import (
	"fmt"
	"sync"
)

var (
	mu1 sync.Mutex
	mu2 sync.Mutex
)

// FIXED version: both goroutines take the locks in the SAME order (mu1 then mu2).
// The deadlock-prone version would take them in opposite orders:
//
//	g1: mu1.Lock(); mu2.Lock()
//	g2: mu2.Lock(); mu1.Lock()  // <- cycle in the Holt diagram -> deadlock
func normalize(a, b *sync.Mutex) {
	a.Lock()
	defer a.Unlock()
	b.Lock()
	defer b.Unlock()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1000)
	for i := 0; i < 500; i++ {
		go func() {
			defer wg.Done()
			normalize(&mu1, &mu2)
		}()
		go func() {
			defer wg.Done()
			normalize(&mu1, &mu2) // same order, no deadlock
		}()
	}
	wg.Wait()
	fmt.Println("done: consistent lock ordering avoids deadlock")
}
