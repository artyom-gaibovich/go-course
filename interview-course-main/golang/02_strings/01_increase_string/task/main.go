// Задача: Какие проблемы в этом коде?

package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	makeString()
	fmt.Println(time.Since(start))
}

func makeString() string {
	str := ""
	for i := 0; i < 100_000; i++ {
		str += fmt.Sprintf("%d", i)
	}
	return str
}
