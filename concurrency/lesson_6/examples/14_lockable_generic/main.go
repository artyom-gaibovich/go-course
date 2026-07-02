package main

import (
	"fmt"
	"sync"
)

// A reusable "synchronized value": here embedding the mutex is fine,
// because lock/unlock IS the public contract of this helper.
type Lockable[T any] struct {
	sync.Mutex
	Value T
}

func main() {
	var counter Lockable[int]
	var wg sync.WaitGroup

	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()
			counter.Lock()
			counter.Value++
			counter.Unlock()
		}()
	}
	wg.Wait()
	fmt.Println("int counter:", counter.Value)

	var label Lockable[string]
	label.Lock()
	label.Value = "synchronized string"
	label.Unlock()
	fmt.Println("string value:", label.Value)
}
