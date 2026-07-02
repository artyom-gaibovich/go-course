package main

import (
	"fmt"
	"net/http"
	"sync"
)

func fetch(wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := http.Get("https://www.google.com")
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
}

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go fetch(&wg)
	}

	wg.Wait()
}
