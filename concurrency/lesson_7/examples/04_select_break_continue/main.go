package main

import "fmt"

func main() {
	ch := make(chan int)
	go func() {
		for i := 1; i <= 4; i++ {
			ch <- i
		}
		close(ch)
	}()

loop:
	for {
		select {
		case v, open := <-ch:
			if !open {
				fmt.Println("channel closed")
				break loop
			}
			switch v {
			case 2:
				continue
			case 3:
				fmt.Println("break only exits this case, not the loop")
				break
			}
			fmt.Println("value:", v)
		}
	}
}