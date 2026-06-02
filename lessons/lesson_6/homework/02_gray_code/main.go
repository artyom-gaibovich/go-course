/*
Домашнее задание 2: Код Грея (Gray Code)

Код Грея — это система счисления, в которой соседние числа отличаются ровно одним битом.
Используется в энкодерах, минимизации булевых функций (карты Карно), исправлении ошибок.

Обычный код:  0 = 000, 1 = 001, 2 = 010, 3 = 011, 4 = 100 ...
Код Грея:     0 = 000, 1 = 001, 2 = 011, 3 = 010, 4 = 110 ...
              (соседи отличаются ровно на 1 бит)

Реализуй:

  func ToGray(n uint) uint
    Конвертирует обычное число в код Грея.
    Формула: gray = n ^ (n >> 1)

  func FromGray(g uint) uint
    Конвертирует код Грея обратно в обычное число.
    Алгоритм: сдвигать g вправо, XOR-ить с предыдущим результатом, пока g != 0.

  func PopCount(a, b uint) int
    Вспомогательная: считает количество отличающихся битов между a и b.
    Подсказка: PopCount(a ^ b) — количество установленных битов в XOR.

Критерии:
  [ ] FromGray(ToGray(n)) == n для всех n в тестах
  [ ] ToGray(n) и ToGray(n+1) отличаются ровно одним битом
      (проверить: diffBits(ToGray(n), ToGray(n+1)) == 1)
  [ ] Каждая функция прокомментирована с объяснением алгоритма

Пример:
  ToGray(0) = 0    (000 → 000)
  ToGray(1) = 1    (001 → 001)
  ToGray(2) = 3    (010 → 011)
  ToGray(3) = 2    (011 → 010)
  ToGray(4) = 6    (100 → 110)
*/

package main

import "fmt"

// ToGray конвертирует число n в код Грея.
// Формула: gray = n ^ (n >> 1)
func ToGray(n uint) uint {
	// TODO
	return 0
}

// FromGray конвертирует код Грея g обратно в обычное число.
func FromGray(g uint) uint {
	// TODO
	return 0
}

// diffBits возвращает количество битов, которые отличаются у a и b.
func diffBits(a, b uint) int {
	// TODO: используй a ^ b и подсчёт единичных битов
	return 0
}

func main() {
	fmt.Println("n\tGray\tBack\tOK\tdiffBits")
	for n := uint(0); n < 16; n++ {
		g := ToGray(n)
		back := FromGray(g)
		ok := back == n

		diff := -1
		if n > 0 {
			diff = diffBits(ToGray(n-1), g)
		}

		fmt.Printf("%d\t%03b\t%d\t%v\t%d\n", n, g, back, ok, diff)
	}
}
