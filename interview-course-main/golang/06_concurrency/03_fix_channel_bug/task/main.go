// Задача: Что выведет программа и почему?

package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan struct{})
	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("continued")
		ch <- struct{}{}
	}()

	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			fmt.Println("tick")
			ch <- struct{}{}
		case value := <-ch:
			fmt.Printf("value %t\n", value)
			return
		}
	}
}
