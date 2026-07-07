package main

import "fmt"

func safe(name string, fn func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("%s -> panic: %v\n", name, r)
		}
	}()
	fn()
	fmt.Printf("%s -> ok\n", name)
}

func main() {
	source := make(chan int)
	clone := source
	fmt.Println("channels are pointers, copy == original:", source == clone)
	fmt.Println("different make != :", source == make(chan int))

	closed := make(chan int)
	close(closed)
	v, ok := <-closed
	fmt.Printf("read from closed: value=%d ok=%v\n", v, ok)

	safe("write to closed", func() { closed <- 1 })
	safe("close twice", func() { close(closed) })

	var nilCh chan int
	safe("close(nil)", func() { close(nilCh) })

	fmt.Println("note: send/recv on nil channel block forever (not shown to avoid deadlock)")
}
