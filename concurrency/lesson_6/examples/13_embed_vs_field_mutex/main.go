package main

import (
	"fmt"
	"sync"
)

// BAD: embedding leaks Lock/Unlock/TryLock to users of Data.
//
//	type Data struct {
//	    sync.Mutex
//	    items []int
//	}
//
// GOOD: named private field, mutex declared ABOVE the field it guards.
type Data struct {
	mu    sync.Mutex // guards items
	items []int
}

func (d *Data) Append(v int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.items = append(d.items, v)
}

// ForEach exposes the guarded slice safely without handing it out.
func (d *Data) ForEach(action func(v int)) {
	if action == nil {
		return
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, v := range d.items {
		action(v)
	}
}

func main() {
	var d Data
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(v int) {
			defer wg.Done()
			d.Append(v)
		}(i)
	}
	wg.Wait()

	sum := 0
	d.ForEach(func(v int) { sum += v })
	fmt.Println("sum:", sum)
}
