package main

import "fmt"

func main() {
	ch := make(chan int, 2)
	ch <- 10
	ch <- 20
	fmt.Printf("len=%d cap=%d (no block, buffer not full)\n", len(ch), cap(ch))

	close(ch)
	for v := range ch {
		fmt.Println(v)
	}
}
