package main

import (
	"fmt"
	"sync"
)

func producer(out chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 5; i++ {
		out <- i
	}
	close(out)
}

func consumer(in <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for v := range in {
		fmt.Println("consumed:", v)
	}
	fmt.Println("range finished: channel closed")
}

func main() {
	ch := make(chan int)
	var wg sync.WaitGroup
	wg.Add(2)
	go producer(ch, &wg)
	go consumer(ch, &wg)
	wg.Wait()
}