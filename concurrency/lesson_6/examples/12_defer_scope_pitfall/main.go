package main

import (
	"fmt"
	"sync"
)

var mu sync.Mutex

// BAD: defer fires at function exit, NOT at block/iteration exit.
//
//	for _, v := range items {
//	    mu.Lock()
//	    defer mu.Unlock() // accumulates, second iteration self-deadlocks
//	    process(v)
//	}
//
// GOOD: scope the lock to one iteration with an explicit Unlock
// (or a withLock wrapper).
func appendAll(items []int) []int {
	var out []int
	for _, v := range items {
		mu.Lock()
		out = append(out, v*v)
		mu.Unlock()
	}
	return out
}

func main() {
	fmt.Println(appendAll([]int{1, 2, 3, 4}))
}
