package main

import (
	"fmt"
	"sync"
)

type LazyMap struct {
	once sync.Once
	mu   sync.Mutex
	data map[string]string
}

func (m *LazyMap) init() {
	m.once.Do(func() {
		fmt.Println("allocating the underlying map (lazily, on first Add)")
		m.data = make(map[string]string)
	})
}

func (m *LazyMap) Add(key, value string) {
	m.init()
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}

func main() {
	var m LazyMap
	fmt.Println("constructed, no allocation yet")
	m.Add("a", "1")
	m.Add("b", "2")
	fmt.Println("size:", len(m.data))
}
