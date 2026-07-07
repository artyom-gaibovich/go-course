package main

import (
	"fmt"
	"time"
)

func fetch(label string, delay time.Duration) <-chan string {
	ch := make(chan string, 1)
	go func() {
		time.Sleep(delay)
		ch <- label
	}()
	return ch
}

func main() {
	start := time.Now()
	c1 := fetch("a", 2*time.Second)
	c2 := fetch("b", 2*time.Second)
	r1, r2 := <-c1, <-c2
	fmt.Printf("concurrent: %s %s in %v\n", r1, r2, time.Since(start).Round(time.Second))

	start = time.Now()
	s1 := <-fetch("a", 2*time.Second)
	s2 := <-fetch("b", 2*time.Second)
	fmt.Printf("sequential: %s %s in %v\n", s1, s2, time.Since(start).Round(time.Second))
}
