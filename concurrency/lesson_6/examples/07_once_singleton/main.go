package main

import (
	"fmt"
	"sync"
)

type singleton struct {
	id int
}

var (
	instance *singleton
	once     sync.Once
	created  int
)

func getInstance() *singleton {
	once.Do(func() {
		created++
		instance = &singleton{id: 42}
	})
	return instance
}

func main() {
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			_ = getInstance()
		}()
	}
	wg.Wait()

	fmt.Println("instance id:", getInstance().id)
	fmt.Println("times created:", created)
}
