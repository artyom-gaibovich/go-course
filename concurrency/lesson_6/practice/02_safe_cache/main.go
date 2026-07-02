/*
Задание 2: Потокобезопасный кэш с правильной инкапсуляцией.

Реализуйте SafeCache (обёртка над map[string]string) с методами:
  - Set(key, value string)
  - Get(key string) (string, bool)
  - Len() int

Требования:
  - Мьютекс — приватное ИМЕНОВАННОЕ поле (не встраивание sync.Mutex).
  - Все методы на УКАЗАТЕЛЬНОМ ресивере.
  - Ленивая инициализация мапы через sync.Once (мапа создаётся при первом Set).
  - Len тоже под блокировкой.
  - Не отдавайте внутреннюю мапу наружу.

Проверьте: go run -race main.go
*/

package main

import (
	"fmt"
	"sync"
)

type SafeCache struct {
	mu   sync.Mutex
	once sync.Once
	data map[string]string
}

func (c *SafeCache) initOnce() {
	// TODO: once.Do(func(){ c.data = make(...) })
}

func (c *SafeCache) Set(key, value string) {
	// TODO
}

func (c *SafeCache) Get(key string) (string, bool) {
	// TODO
	return "", false
}

func (c *SafeCache) Len() int {
	// TODO
	return 0
}

func main() {
	var c SafeCache
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func(n int) {
			defer wg.Done()
			c.Set(fmt.Sprintf("key-%d", n), "v")
		}(i)
	}
	wg.Wait()
	fmt.Println("len:", c.Len())
}
