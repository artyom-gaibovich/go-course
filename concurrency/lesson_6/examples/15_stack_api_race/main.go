package main

import (
	"fmt"
	"sync"
)

// Each method is thread-safe, but separate Top()+Pop() calls are NOT atomic
// together: another goroutine can slip between them. The fix is a single
// atomic Pop that both removes and returns the value under one lock.
type Stack struct {
	mu   sync.Mutex
	data []int
}

func (s *Stack) Push(v int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = append(s.data, v)
}

func (s *Stack) Pop() (int, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.data) == 0 {
		return 0, false
	}
	v := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return v, true
}

func main() {
	var s Stack
	for i := 0; i < 1000; i++ {
		s.Push(i)
	}

	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				s.Pop()
			}
		}()
	}
	wg.Wait()

	_, ok := s.Pop()
	fmt.Println("remaining items, pop ok:", ok)
}
