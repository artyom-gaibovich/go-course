/*
Задание 2: Статистика через variadic

Реализуй пакет статистических функций с variadic-параметром:

  sum(nums ...float64) float64          — сумма всех чисел
  average(nums ...float64) float64      — среднее арифметическое (0 если нет аргументов)
  min(nums ...float64) float64          — минимум
  max(nums ...float64) float64          — максимум

Для min/max: если аргументов нет — паникуй с сообщением "нет аргументов".

Затем напиши функцию stats(nums ...float64), которая выводит все четыре значения сразу.

Пример вывода для stats(3, 1, 4, 1, 5, 9, 2, 6):
  Сумма:    31.00
  Среднее:  3.88
  Минимум:  1.00
  Максимум: 9.00

Подсказка: stats может передавать nums... во вложенные функции через spread.
*/

package main

import "fmt"

func sum(nums ...float64) float64 {
	total := 0.0
	for _, n := range nums {
		total += n
	}
	return total
}

func average(nums ...float64) float64 {
	if len(nums) == 0 {
		return 0
	}
	// TODO: используй sum
	return 0
}

func min(nums ...float64) float64 {
	if len(nums) == 0 {
		panic("нет аргументов")
	}
	// TODO
	return 0
}

func max(nums ...float64) float64 {
	if len(nums) == 0 {
		panic("нет аргументов")
	}
	// TODO
	return 0
}

func stats(nums ...float64) {
	// TODO: выведи sum, average, min, max — передавай nums... в каждую функцию
	fmt.Println("Сумма:   ", sum(nums...))
}

func main() {
	stats(3, 1, 4, 1, 5, 9, 2, 6)
	fmt.Println()

	nums := []float64{10, 20, 30, 40, 50}
	stats(nums...) // spread слайса
}
