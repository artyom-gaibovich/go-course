package main

import (
	"fmt"
	"sync"
)

// BAD: receiving the WaitGroup by value copies the counter.
// Done() decrements the COPY; the original never reaches zero -> Wait() hangs.
//
//	func worker(wg sync.WaitGroup) { defer wg.Done(); ... } // deadlock!
//
// GOOD: pass by pointer so every goroutine touches the same counter.
func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("worker %d done\n", id)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go worker(i, &wg)
	}
	wg.Wait()
	fmt.Println("finished: WaitGroup passed by pointer, never copied")
}
