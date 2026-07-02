package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		fmt.Println("1-горутина")
		wg.Done()
	}()

	go func() {
		fmt.Println("2-горутина")
		wg.Done()
	}()

	wg.Wait()
	fmt.Println("all workers finished")
}
