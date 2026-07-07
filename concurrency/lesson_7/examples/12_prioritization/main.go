package main

import "fmt"

func main() {
	high := make(chan int, 10)
	low := make(chan int, 10)
	for i := 0; i < 5; i++ {
		high <- 100 + i
		low <- i
	}

	drained := 0
	for drained < 10 {
		select {
		case v := <-high:
			fmt.Println("HIGH:", v)
			drained++
		default:
			select {
			case v := <-high:
				fmt.Println("HIGH:", v)
				drained++
			case v := <-low:
				fmt.Println("low:", v)
				drained++
			}
		}
	}
}