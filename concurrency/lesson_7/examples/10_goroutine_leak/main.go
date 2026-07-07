package main

import (
	"fmt"
	"runtime"
	"time"
)

func leaky() {
	ch := make(chan string)
	go func() {
		for v := range ch {
			fmt.Println("worker:", v)
		}
		fmt.Println("worker exited")
	}()
	ch <- "hello"
}

func fixed() {
	ch := make(chan string)
	go func() {
		for v := range ch {
			fmt.Println("worker:", v)
		}
		fmt.Println("worker exited")
	}()
	ch <- "hello"
	close(ch)
}

func main() {
	leaky()
	time.Sleep(50 * time.Millisecond)
	fmt.Println("after leaky, goroutines:", runtime.NumGoroutine())

	fixed()
	time.Sleep(50 * time.Millisecond)
	fmt.Println("after fixed, goroutines:", runtime.NumGoroutine())
}
