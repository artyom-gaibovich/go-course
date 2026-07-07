package main

import "fmt"

func produce(out chan<- int) {
	for i := 0; i < 3; i++ {
		out <- i
	}
	close(out)
}

func consume(in <-chan int, done chan<- struct{}) {
	for v := range in {
		fmt.Println("received:", v)
	}
	close(done)
}

func main() {
	ch := make(chan int)
	done := make(chan struct{})

	go produce(ch)
	go consume(ch, done)

	<-done
	fmt.Println("all consumed")
}
