// Задача: Что выведет программа и почему? Откуда возникает рандом?

package main

import "fmt"

func main() {
	myMap := map[int]string{1: "a", 2: "b", 3: "c"}

	for key, value := range myMap {
		fmt.Println("Key:", key, "Value:", value)
	}
}
