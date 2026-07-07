/*
Задание 5: First-response-wins без утечек горутин.

Запустите 5 горутин-«реплик», каждая с разной задержкой возвращает свой ответ.
Верните первый пришедший ответ и НЕ допустите утечки остальных четырёх горутин:
проигравшие не должны навсегда зависнуть на записи в канал.

В комментарии объясните, почему выбран буферизированный канал ёмкостью на всех
писателей (а не небуферизированный и не select с default).

Подсказка: make(chan T, 5) — проигравшие положат результат в буфер, не
блокируясь, завершатся; невостребованный буфер потом соберёт GC.
*/

package main

import (
	"fmt"
	"math/rand"
	"time"
)

func firstResponse() string {
	ch := make(chan string)

	for i := 0; i < 5; i++ {
		go func(id int) {
			randomSeconds := rand.Intn(5) + 1 // 1-5
			duration := time.Duration(randomSeconds) * time.Second
			time.Sleep(duration)
			ch <- fmt.Sprintf("Горутина %d спала %d секунд", id, randomSeconds)
		}(i)
	}

	v := <-ch
	return v
	// TODO: буферизированный канал на 5 значений
	// TODO: запустить 5 горутин с разной задержкой, каждая пишет свой ответ
	// TODO: вернуть первое прочитанное значение
}

func main() {
	// TODO: вызвать firstResponse и напечатать результат
	fmt.Println(firstResponse())
}
