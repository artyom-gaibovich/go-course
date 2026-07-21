package main

import "fmt"

// Задача: Оптимальное планирование отпуска
//
// Необходимо определить день X начала отпуска длиной в V дней так, чтобы отгулять весь отпуск
// в ближайшие P дней и пропустить минимум Y встреч.
// Считаем, что уже завтра первый возможный день отпуска (day: 1).
//
// Входные данные:
// - daysWithMeetings: []DayMeetings - дни со встречами (уже упорядочены по дню)
//   - day: int - номер дня
//   - meetings: int - количество встреч в этот день
// - periodLength: int - в какой срок надо отгулять ВЕСЬ отпуск (в ближайшие P дней)
// - vacationLength: int - длительность отпуска в днях
//
// Выходные данные:
// - []int - [день X начала отпуска, сколько встреч Y пропустим]
//
// Примеры:
// daysWithMeetings = [{day: 3, meetings: 1}, {day: 4, meetings: 3}, {day: 14, meetings: 3}, {day: 21, meetings: 3}, {day: 28, meetings: 1}]
// periodLength = 30
// vacationLength = 7
// Результат: [5, 0] - начать отпуск с 5 дня, пропустить 0 встреч

type DayMeetings struct {
	Day      int
	Meetings int
}

func findOptimalVacation(daysWithMeetings []DayMeetings, periodLength int, vacationLength int) []int {
	// TODO: реализуйте функцию
	return nil
}

func main() {
	// Пример 1
	daysWithMeetings1 := []DayMeetings{
		{Day: 3, Meetings: 1},
		{Day: 4, Meetings: 3},
		{Day: 14, Meetings: 3},
		{Day: 21, Meetings: 3},
		{Day: 28, Meetings: 1},
	}
	result1 := findOptimalVacation(daysWithMeetings1, 30, 7)
	// Ожидается: [5, 0]
	fmt.Println("Test 1:", result1)

	// Пример 2
	daysWithMeetings2 := []DayMeetings{
		{Day: 3, Meetings: 1},
		{Day: 4, Meetings: 3},
		{Day: 5, Meetings: 3},
		{Day: 9, Meetings: 5},
		{Day: 13, Meetings: 2},
		{Day: 14, Meetings: 1},
		{Day: 21, Meetings: 3},
		{Day: 25, Meetings: 3},
		{Day: 28, Meetings: 6},
	}
	result2 := findOptimalVacation(daysWithMeetings2, 31, 14)
	// Ожидается: [10, 6]
	fmt.Println("Test 2:", result2)
}
