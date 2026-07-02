/*
Домашнее задание 2: Ленивый синглтон.

Сделайте конфиг-загрузчик GetConfig(), который читает "тяжёлый" конфиг
ровно ОДИН РАЗ через sync.Once, даже если его дёргают из сотни горутин
одновременно.

Требования:
  - Инициализация выполняется лениво (при первом GetConfig).
  - Используйте sync.Once.
  - Докажите однократность: счётчик инициализаций должен быть == 1.

Проверьте: go run -race main.go
*/

package main

import (
	"fmt"
	"sync"
)

type Config struct {
	Endpoint string
}

var (
	cfg     *Config
	once    sync.Once
	initCnt int
)

func loadConfig() *Config {
	// TODO: once.Do(func(){ initCnt++; cfg = &Config{...} })
	return cfg
}

func GetConfig() *Config {
	return loadConfig()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			_ = GetConfig()
		}()
	}
	wg.Wait()

	fmt.Println("endpoint:", GetConfig().Endpoint)
	fmt.Println("init count (must be 1):", initCnt)
}
