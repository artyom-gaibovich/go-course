package main

import "fmt"

func tryRecv(ch <-chan int) (int, bool) {
	select {
	case v := <-ch:
		return v, true
	default:
		return 0, false
	}
}

func trySend(ch chan<- int, v int) bool {
	select {
	case ch <- v:
		return true
	default:
		return false
	}
}

func main() {
	ch := make(chan int, 1)

	fmt.Println("trySend 1:", trySend(ch, 1))
	fmt.Println("trySend 2 (full):", trySend(ch, 2))

	v, ok := tryRecv(ch)
	fmt.Printf("tryRecv: value=%d ok=%v\n", v, ok)

	v, ok = tryRecv(ch)
	fmt.Printf("tryRecv (empty): value=%d ok=%v\n", v, ok)
}