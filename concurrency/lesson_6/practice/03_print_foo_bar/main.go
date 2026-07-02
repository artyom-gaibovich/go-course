/*
Задание 3: Поочерёдная печать FOO/BAR.

Запустите две горутины: одна печатает "FOO", другая "BAR".
Нужно, чтобы вывод чередовался строго: FOO BAR FOO BAR ... ровно N раз.

Подсказка: чередование — это race condition ПО ДИЗАЙНУ. Нужно явно
передавать "ход" от одной горутины к другой. Можно двумя мьютексами
(один изначально захвачен) или через WaitGroup по шагам.

Ожидаемый вывод (N=3): FOO BAR FOO BAR FOO BAR
*/

package main

import (
	"sync"
)

func printFooBar(n int) {
	var fooLock, barLock sync.Mutex
	barLock.Lock() // BAR ждёт, пока FOO передаст ход

	// TODO: горутина FOO — печатает "FOO", затем разблокирует barLock
	// TODO: горутина BAR — ждёт barLock, печатает "BAR", разблокирует fooLock
	// TODO: повторить n раз, дождаться завершения
	_ = fooLock
}

func main() {
	printFooBar(3)
}
