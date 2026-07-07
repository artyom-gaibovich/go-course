package main

import (
	"fmt"
	"sync"
)

func main() {
	signal := make(chan struct{})
	var wg sync.WaitGroup

	for id := 1; id <= 3; id++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			<-signal
			fmt.Printf("subscriber %d notified\n", id)
		}(id)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		close(signal)
	}()

	wg.Wait()
}
