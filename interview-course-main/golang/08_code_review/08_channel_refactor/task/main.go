// Задача: Что выведет программа и почему? Отрефакторить код.

package main

import (
	"fmt"
	"strconv"
	"sync"
)

func main() {
	var wc sync.WaitGroup
	fff := sync.Mutex{}
	a := make(chan string, 3)

	for i := 0; i < 5; i++ {
		wc.Add(1)
		go func(a chan<- string, i int, wc *sync.WaitGroup) {
			defer wc.Done()
			fff.Lock()
			a <- fmt.Sprintf("Current gorutine number: %s", strconv.Itoa(i))
			fff.Unlock()
		}(a, i, &wc)
	}

	for {
		select {
		case result := <-a:
			fmt.Printf(result)
		}
	}

	wc.Wait()
}
