package main

import (
	"fmt"
	"sync"
)

// sync.Mutex is NOT recursive. Get() locks, then calls Size().
// If Size() locked again -> self-deadlock. The fix: a private,
// lock-free sizeLocked() called while the lock is already held.
type Cache struct {
	mu   sync.Mutex
	data map[string]string
}

func NewCache() *Cache {
	return &Cache{data: make(map[string]string)}
}

func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

// sizeLocked must be called with c.mu already held.
func (c *Cache) sizeLocked() int {
	return len(c.data)
}

func (c *Cache) Size() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.sizeLocked()
}

func (c *Cache) Get(key string) string {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.sizeLocked() == 0 {
		return ""
	}
	return c.data[key]
}

func main() {
	c := NewCache()
	c.Set("a", "1")
	c.Set("b", "2")
	fmt.Println("size:", c.Size())
	fmt.Println("get a:", c.Get("a"))
}
