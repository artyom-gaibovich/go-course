/*
Задание 2: Fan-in нескольких источников через select.

Есть три горутины-генератора, каждая пишет несколько чисел в свой канал и
закрывает его. Соберите все значения в main через select в цикле и печатайте их.
Корректно завершитесь, когда все три канала закрыты.

Подсказка: «выключайте» закрытый канал, присваивая его переменную nil — чтение
из nil-канала блокируется навсегда, поэтому такой кейс перестаёт выбираться.
Считайте число ещё открытых каналов и выходите, когда оно станет 0.
*/

package main

import "fmt"

func generator(values ...int) <-chan int {
	ch := make(chan int, len(values))
	// TODO: записать values в ch и закрыть канал
	for _, value := range values {
		ch <- value
	}
	close(ch)
	return ch
}

func main() {
	ch1 := generator(1, 2)
	ch2 := generator(3, 4)
	ch3 := generator(5, 6)

	for {
		select {
		case v, ok := <-ch1:
			if !ok {
				fmt.Println("ch1 closed")
				ch1 = nil
				continue
			}

			fmt.Println(v)
		case v, ok := <-ch2:
			if !ok {
				fmt.Println("ch2 closed")
				ch2 = nil
				continue
			}
			fmt.Println(v)
		case v, ok := <-ch3:
			if !ok {
				fmt.Println("ch3 closed")
				ch3 = nil
				continue
			}
			fmt.Println(v)
		default:
			return
		}
	}

	// TODO: создать три генератора
	// TODO: цикл с select, читать из всех каналов, закрытый канал = nil
	// TODO: завершиться, когда все каналы закрыты
}
