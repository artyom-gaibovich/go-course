package main

import "fmt"

// ============================================================================
// Решение: Подсчет участников и суммирование шагов
// ============================================================================

func findChampions(statistics [][]UserSteps) Champions {
	// Подсчитываем количество дней соревнований
	totalDays := len(statistics)

	// Подсчитываем, в скольких днях участвовал каждый пользователь
	userDaysCount := make(map[int]int)
	userTotalSteps := make(map[int]int)

	// Проходим по всем дням и собираем статистику
	for _, dayStats := range statistics {
		// Создаем set участников этого дня
		dayParticipants := make(map[int]bool)
		for _, userSteps := range dayStats {
			dayParticipants[userSteps.UserId] = true
			userTotalSteps[userSteps.UserId] += userSteps.Steps
		}

		// Увеличиваем счетчик дней для всех участников этого дня
		for userId := range dayParticipants {
			userDaysCount[userId]++
		}
	}

	// Находим пользователей, участвовавших во всех днях
	var candidates []int
	for userId, daysCount := range userDaysCount {
		if daysCount == totalDays {
			candidates = append(candidates, userId)
		}
	}

	// Если нет кандидатов, возвращаем пустой результат
	if len(candidates) == 0 {
		return Champions{UserIds: []int{}, Steps: 0}
	}

	// Находим максимальную сумму шагов среди кандидатов
	maxSteps := 0
	for _, userId := range candidates {
		if userTotalSteps[userId] > maxSteps {
			maxSteps = userTotalSteps[userId]
		}
	}

	// Находим всех чемпионов с максимальной суммой шагов
	var champions []int
	for _, userId := range candidates {
		if userTotalSteps[userId] == maxSteps {
			champions = append(champions, userId)
		}
	}

	return Champions{
		UserIds: champions,
		Steps:   maxSteps,
	}
}

// ============================================================================
// Демонстрация и тесты
// ============================================================================

func main() {
	fmt.Println("=== Решение: Определение чемпионов ===")
	testFindChampions()
}

func testFindChampions() {
	// Тест 1
	statistics1 := [][]UserSteps{
		{{UserId: 1, Steps: 1000}, {UserId: 2, Steps: 1500}},
		{{UserId: 2, Steps: 1000}},
	}
	result1 := findChampions(statistics1)
	fmt.Printf("Тест 1: UserIds=%v, Steps=%d (ожидается [2], 2500)\n", result1.UserIds, result1.Steps)

	// Тест 2
	statistics2 := [][]UserSteps{
		{{UserId: 1, Steps: 2000}, {UserId: 2, Steps: 1500}},
		{{UserId: 2, Steps: 4000}, {UserId: 1, Steps: 3500}},
	}
	result2 := findChampions(statistics2)
	fmt.Printf("Тест 2: UserIds=%v, Steps=%d (ожидается [1, 2], 5500)\n", result2.UserIds, result2.Steps)
}

// ============================================================================
// Объяснение решения
// ============================================================================

// РЕШЕНИЕ: Подсчет участников и суммирование шагов
//
// Идея:
// - Нужно найти пользователей, которые участвовали во ВСЕХ днях
// - Среди таких пользователей выбрать тех, у кого максимальная сумма шагов
// - Используем два map: для подсчета дней участия и для суммирования шагов
//
// Алгоритм:
// 1. Подсчитываем общее количество дней соревнований
// 2. Для каждого дня:
//    - Собираем участников этого дня в set
//    - Увеличиваем счетчик дней для каждого участника
//    - Суммируем шаги каждого участника
// 3. Находим кандидатов (участников всех дней)
// 4. Находим максимальную сумму шагов среди кандидатов
// 5. Возвращаем всех кандидатов с максимальной суммой шагов
//
// Пример работы:
// statistics = [
//   [{userId: 1, steps: 1000}, {userId: 2, steps: 1500}],  // день 1
//   [{userId: 2, steps: 1000}]                              // день 2
// ]
//
// Шаг 1: Обрабатываем день 1
//   dayParticipants = {1: true, 2: true}
//   userDaysCount[1] = 1
//   userDaysCount[2] = 1
//   userTotalSteps[1] = 1000
//   userTotalSteps[2] = 1500
//
// Шаг 2: Обрабатываем день 2
//   dayParticipants = {2: true}
//   userDaysCount[2] = 2  (участвовал в обоих днях!)
//   userTotalSteps[2] = 1500 + 1000 = 2500
//
// Шаг 3: Находим кандидатов
//   totalDays = 2
//   userDaysCount[1] = 1 ≠ 2 → не кандидат
//   userDaysCount[2] = 2 = 2 → кандидат!
//   candidates = [2]
//
// Шаг 4: Находим максимальную сумму шагов
//   maxSteps = userTotalSteps[2] = 2500
//
// Шаг 5: Находим чемпионов
//   champions = [2] (единственный кандидат с maxSteps)
//
// Результат: {userIds: [2], steps: 2500}
//
// Временная сложность: O(d * u)
// - d - количество дней
// - u - среднее количество участников в день
// - Для каждого дня проходим по участникам: O(d * u)
// - Поиск кандидатов и чемпионов: O(u)
// - Итого: O(d * u)
//
// Пространственная сложность: O(u)
// - userDaysCount: O(u)
// - userTotalSteps: O(u)
// - dayParticipants: O(u) для каждого дня, но переиспользуется
// - Итого: O(u)
//
// Почему используем set для dayParticipants?
// - В одном дне может быть несколько записей для одного пользователя
// - Set гарантирует, что каждый пользователь учитывается только один раз за день
// - Это важно для корректного подсчета дней участия
//
// Пример с дубликатами:
// statistics = [
//   [{userId: 1, steps: 1000}, {userId: 1, steps: 500}],  // день 1
//   [{userId: 1, steps: 2000}]                              // день 2
// ]
// Без set: userDaysCount[1] = 3 (неправильно!)
// С set: userDaysCount[1] = 2 (правильно!)
