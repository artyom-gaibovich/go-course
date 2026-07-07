package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string)
	go func() {
		time.Sleep(2 * time.Second)
		ch <- "some str"
	}()
	value := <-ch
	fmt.Println(value)

}
