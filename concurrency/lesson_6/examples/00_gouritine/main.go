package main

import (
	"fmt"
	"time"
)

func printNames(v int, label string) {
	fmt.Println(v, label)
}

func cycleRun(label string) {
	for i := 0; i < 1000; i++ {
		go printNames(i, label)
	}
}

func main() {
	go cycleRun("first cycle")
	go cycleRun("second cycle")
	go cycleRun("third cycle")

	time.Sleep(5 * time.Second)
}
