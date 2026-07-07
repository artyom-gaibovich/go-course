package main

import (
	"fmt"
	"sync"
)

var mu sync.Mutex
var counter int

func main() {

	var wg sync.WaitGroup
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}

	wg.Wait()
	fmt.Println("counter:", counter)

}
