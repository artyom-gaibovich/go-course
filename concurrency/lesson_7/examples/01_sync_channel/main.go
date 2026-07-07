package main

import (
	"fmt"
	"time"
)

func worker(out chan<- string) {
	time.Sleep(2 * time.Second)
	out <- "job done"
}

func main() {
	results := make(chan string)
	go worker(results)

	fmt.Println("waiting for result...")
	result := <-results
	fmt.Println("got:", result)
}
