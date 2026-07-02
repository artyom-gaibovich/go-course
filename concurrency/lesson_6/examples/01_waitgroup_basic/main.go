package main

import (
	"fmt"
	"sync"
)

func main() {
	const workers = 5

	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func(id int) {
			defer wg.Done()
			fmt.Printf("worker %d done\n", id)
		}(i)
	}

	wg.Wait()
	fmt.Println("all workers finished")
}
