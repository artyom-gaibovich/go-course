// Задача: Напишите функцию GetOrCreate, которая создает новый элемент мапы,
// если его еще не было и возвращает его значение, или просто возвращает значение при наличии.
// Важно учесть, что код должен нормально работать в конкурентной среде.

package main

import (
	"fmt"
	"sync"
)

func main() {
	cm := NewConcurrentMap()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		val := cm.GetOrCreate("key1", "value1")
		fmt.Println("Goroutine 1 got:", val)
	}()

	go func() {
		defer wg.Done()
		val := cm.GetOrCreate("key1", "value2")
		fmt.Println("Goroutine 2 got:", val)
	}()

	wg.Wait()
}
