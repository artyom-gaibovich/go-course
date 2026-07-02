package main

import (
	"fmt"
	"sync"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered panic:", r)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Done()
	wg.Done()

	fmt.Println("unreachable: second Done drove the counter negative")
}
