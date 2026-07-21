// Задача: Провести рефакторинг кода
//
// ТЗ: программа запускает 5 горутин, каждая отправляет сообщение в канал.
// Главная горутина должна прочитать все сообщения и завершить работу.
//
// Этот код НАМЕРЕННО содержит ошибки для учебных целей!
// Не запускайте в production!

package main

import (
	"fmt"
	"strconv"
	"sync"
)

func main() {
	var wc sync.WaitGroup
	m := make(chan string, 3)
	fff := sync.Mutex{}

	for i := 0; i < 5; i++ {
		wc.Add(1)
		go func(mm chan<- string, i int, group *sync.WaitGroup) {
			defer wc.Done()
			fff.Lock()
			mm <- fmt.Sprintf("Gorutine %s", strconv.Itoa(i))
			fff.Unlock()
		}(m, i, &wc)
	}

	for {
		select {
		case q := <-m:
			fmt.Println(q)
		}
	}

	wc.Wait()
	close(m)
}
