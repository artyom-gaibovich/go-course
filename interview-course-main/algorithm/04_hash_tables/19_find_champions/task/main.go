package main

import "fmt"

// Задача: Определение чемпионов по шагам
//
// Необходимо определить userids участников, которые прошли наибольшее количество шагов
// за все дни, не пропустив ни одного дня соревнований.
//
// Входные данные:
// - statistics: [][]UserSteps - статистика по дням
//   - каждый внутренний массив содержит статистику за один день
//   - UserSteps содержит userId и steps (количество шагов за день)
//
// Выходные данные:
// - Champions структура:
//   - UserIds: []int - список ID пользователей-чемпионов
//   - Steps: int - общее количество шагов чемпионов
//
// Требования:
// - Чемпионом может стать только тот, кто участвовал во ВСЕХ днях соревнований
// - Среди таких участников выбираются те, у кого максимальная сумма шагов
//
// Примеры:
// statistics = [[{userId: 1, steps: 1000}, {userId: 2, steps: 1500}], [{userId: 2, steps: 1000}]]
// Результат: {userIds: [2], steps: 2500}
//
// statistics = [[{userId: 1, steps: 2000}, {userId: 2, steps: 1500}], [{userId: 2, steps: 4000}, {userId: 1, steps: 3500}]]
// Результат: {userIds: [1, 2], steps: 5500}

type UserSteps struct {
	UserId int
	Steps  int
}

type Champions struct {
	UserIds []int
	Steps   int
}

func findChampions(statistics [][]UserSteps) Champions {
	// TODO: реализуйте функцию
	return Champions{}
}

func main() {
	// Пример 1
	statistics1 := [][]UserSteps{
		{{UserId: 1, Steps: 1000}, {UserId: 2, Steps: 1500}},
		{{UserId: 2, Steps: 1000}},
	}
	result1 := findChampions(statistics1)
	// Ожидается: {userIds: [2], steps: 2500}
	fmt.Printf("Test 1: UserIds=%v, Steps=%d\n", result1.UserIds, result1.Steps)

	// Пример 2
	statistics2 := [][]UserSteps{
		{{UserId: 1, Steps: 2000}, {UserId: 2, Steps: 1500}},
		{{UserId: 2, Steps: 4000}, {UserId: 1, Steps: 3500}},
	}
	result2 := findChampions(statistics2)
	// Ожидается: {userIds: [1, 2], steps: 5500}
	fmt.Printf("Test 2: UserIds=%v, Steps=%d\n", result2.UserIds, result2.Steps)
}
