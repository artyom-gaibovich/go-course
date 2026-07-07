package main

import (
	"fmt"
	"sync"
)

func withMutex() int {
	var mu sync.Mutex
	var wg sync.WaitGroup
	value := 0
	wg.Add(2)
	for i := 0; i < 2; i++ {
		go func() {
			defer wg.Done()
			mu.Lock()
			value++
			mu.Unlock()
		}()
	}
	wg.Wait()
	return value
}

func withChannel() int {
	ch := make(chan int)
	go func() { ch <- 1 }()
	go func() { ch <- 1 }()

	value := 0
	value += <-ch
	value += <-ch
	return value
}

func main() {
	fmt.Println("shared memory + mutex:", withMutex())
	fmt.Println("share by communicating:", withChannel())
}
