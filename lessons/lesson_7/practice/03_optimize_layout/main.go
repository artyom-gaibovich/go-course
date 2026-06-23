/*
Задание 3. Оптимизация раскладки структуры (выравнивание).

Дана структура Record. Из-за порядка полей компилятор вставляет padding,
и размер получается больше необходимого.

Задача: создай структуру OptimizedRecord с ТЕМИ ЖЕ полями, но переставь их так,
чтобы unsafe.Sizeof(OptimizedRecord{}) был минимальным.
Программа должна напечатать оба размера; оптимизированный обязан быть меньше.

Подсказка: располагай поля от «широких» к «узким» (int64 → int32 → int16 → bool).
Каждое поле выравнивается по своему размеру; группировка одинаковых ширин
убирает дыры выравнивания.
*/

package main

import (
	"fmt"
	"unsafe"
)

type Record struct {
	flagA bool
	id    int64
	count int32
	flagB bool
	code  int16
}

// TODO: объяви OptimizedRecord с теми же полями в другом порядке
type OptimizedRecord struct {
}

func main() {
	fmt.Println("Record:", unsafe.Sizeof(Record{}))
	fmt.Println("OptimizedRecord:", unsafe.Sizeof(OptimizedRecord{}))
}
