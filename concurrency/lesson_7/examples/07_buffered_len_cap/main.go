package main

import "fmt"

func main() {
	ch := make(chan int, 2)
	ch <- 10
	ch <- 20
	fmt.Printf("len=%d cap=%d (no block, buffer not full)\n", len(ch), cap(ch))

	fmt.Println("read:", <-ch)

	close(ch)
	for v, ok := <-ch; ok; v, ok = <-ch {
		fmt.Println("drained after close:", v)
	}

	v, ok := <-ch
	fmt.Printf("empty closed channel: value=%d ok=%v\n", v, ok)
}
