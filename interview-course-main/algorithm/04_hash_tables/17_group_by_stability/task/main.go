package main

import "fmt"

// Задача: Группировка серверов по стабильности
//
// Напишите функцию, которая по списку метрик серверов строит отображение
// стабильность → список ID серверов.
//
// Входные данные:
// - stats: список структур с полями server (ID сервера) и stability (стабильность)
//
// Выходные данные:
// - map[string][]int: ключ - стабильность (строка), значение - список ID серверов
//
// Примеры:
// stats = [{server: 1, stability: 99}, {server: 2, stability: 97}, {server: 3, stability: 34}, {server: 4, stability: 97}, {server: 5, stability: 97.1}]
// Результат: {"34": [3], "97": [2, 4], "99": [1], "97.1": [5]}

type ServerStat struct {
	Server    int
	Stability float64
}

func groupByStability(stats []ServerStat) map[string][]int {
	// TODO: реализуйте функцию
	return nil
}

func main() {
	stats := []ServerStat{
		{Server: 1, Stability: 99},
		{Server: 2, Stability: 97},
		{Server: 3, Stability: 34},
		{Server: 4, Stability: 97},
		{Server: 5, Stability: 97.1},
	}

	result := groupByStability(stats)
	// Ожидается: {"34": [3], "97": [2, 4], "99": [1], "97.1": [5]}
	fmt.Println("Result:", result)
}
