package main

import (
	"fmt"
	"time"
)

func source(name string, delay time.Duration) <-chan string {
	ch := make(chan string)
	go func() {
		time.Sleep(delay)
		ch <- name
	}()
	return ch
}

func main() {
	c1 := source("fast", 100*time.Millisecond)
	c2 := source("slow", 500*time.Millisecond)

	select {
	case v := <-c1:
		fmt.Println("from c1:", v)
	case v := <-c2:
		fmt.Println("from c2:", v)
	default:
		fmt.Println("nothing ready yet")
	}

	select {
	case v := <-c1:
		fmt.Println("winner:", v)
	case v := <-c2:
		fmt.Println("winner:", v)
	}
}
