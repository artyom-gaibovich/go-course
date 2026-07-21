package main

import "fmt"

// ============================================================================
// Решение: Скользящее окно
// ============================================================================

func findOptimalVacation(daysWithMeetings []DayMeetings, periodLength int, vacationLength int) []int {
	// Создаем массив встреч для каждого дня в периоде
	meetingsByDay := make([]int, periodLength+1) // +1 для индексации с 1
	for _, dm := range daysWithMeetings {
		if dm.Day <= periodLength {
			meetingsByDay[dm.Day] = dm.Meetings
		}
	}

	// Вычисляем сумму встреч для первого окна [1, vacationLength]
	currentMeetings := 0
	for day := 1; day <= vacationLength; day++ {
		currentMeetings += meetingsByDay[day]
	}

	bestStartDay := 1
	minMeetings := currentMeetings

	// Скользящее окно: сдвигаем окно на один день вправо
	for startDay := 2; startDay <= periodLength-vacationLength+1; startDay++ {
		// Убираем встречи из левого конца окна
		currentMeetings -= meetingsByDay[startDay-1]
		// Добавляем встречи из правого конца окна
		currentMeetings += meetingsByDay[startDay+vacationLength-1]

		// Обновляем лучший вариант
		if currentMeetings < minMeetings {
			minMeetings = currentMeetings
			bestStartDay = startDay
		}
	}

	return []int{bestStartDay, minMeetings}
}

// ============================================================================
// Демонстрация и тесты
// ============================================================================

func main() {
	fmt.Println("=== Решение: Оптимальное планирование отпуска ===")
	testFindOptimalVacation()
}

func testFindOptimalVacation() {
	// Тест 1
	daysWithMeetings1 := []DayMeetings{
		{Day: 3, Meetings: 1},
		{Day: 4, Meetings: 3},
		{Day: 14, Meetings: 3},
		{Day: 21, Meetings: 3},
		{Day: 28, Meetings: 1},
	}
	result1 := findOptimalVacation(daysWithMeetings1, 30, 7)
	fmt.Printf("Тест 1: %v (ожидается [5, 0])\n", result1)

	// Тест 2
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
	fmt.Printf("Тест 2: %v (ожидается [10, 6])\n", result2)
}

// ============================================================================
// Объяснение решения
// ============================================================================

// РЕШЕНИЕ: Скользящее окно (Sliding Window)
//
// Идея:
// - Нужно найти окно длиной vacationLength дней с минимальным количеством встреч
// - Используем технику скользящего окна для эффективного перебора всех возможных окон
// - Вместо пересчета суммы встреч для каждого окна, обновляем сумму при сдвиге окна
//
// Алгоритм:
// 1. Создаем массив meetingsByDay для быстрого доступа к количеству встреч по дню
// 2. Вычисляем сумму встреч для первого окна [1, vacationLength]
// 3. Сдвигаем окно на один день вправо:
//    - Убираем встречи из левого конца (startDay-1)
//    - Добавляем встречи из правого конца (startDay+vacationLength-1)
// 4. Обновляем лучший вариант при необходимости
// 5. Возвращаем день начала и количество пропущенных встреч
//
// Пример работы:
// daysWithMeetings = [
//   {day: 3, meetings: 1},
//   {day: 4, meetings: 3},
//   {day: 14, meetings: 3},
//   {day: 21, meetings: 3},
//   {day: 28, meetings: 1}
// ]
// periodLength = 30
// vacationLength = 7
//
// Шаг 1: Создаем массив встреч
//   meetingsByDay = [0, 0, 0, 1, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, ...]
//                   день: 0  1  2  3  4  5  6  7  8  9 10 11 12 13 14
//
// Шаг 2: Первое окно [1, 7]
//   currentMeetings = meetingsByDay[1] + ... + meetingsByDay[7]
//                   = 0 + 0 + 1 + 3 + 0 + 0 + 0 = 4
//   bestStartDay = 1, minMeetings = 4
//
// Шаг 3: Сдвигаем окно на [2, 8]
//   Убираем: meetingsByDay[1] = 0
//   Добавляем: meetingsByDay[8] = 0
//   currentMeetings = 4 - 0 + 0 = 4
//
// Шаг 4: Сдвигаем окно на [3, 9]
//   Убираем: meetingsByDay[2] = 0
//   Добавляем: meetingsByDay[9] = 0
//   currentMeetings = 4 - 0 + 0 = 4
//
// Шаг 5: Сдвигаем окно на [4, 10]
//   Убираем: meetingsByDay[3] = 1
//   Добавляем: meetingsByDay[10] = 0
//   currentMeetings = 4 - 1 + 0 = 3
//
// Шаг 6: Сдвигаем окно на [5, 11]
//   Убираем: meetingsByDay[4] = 3
//   Добавляем: meetingsByDay[11] = 0
//   currentMeetings = 3 - 3 + 0 = 0 ← минимум!
//   bestStartDay = 5, minMeetings = 0
//
// ... продолжаем до окна [24, 30]
//
// Результат: [5, 0] - начать отпуск с 5 дня, пропустить 0 встреч
//
// Временная сложность: O(n + p)
// - n - количество дней со встречами
// - p - periodLength
// - Создание массива meetingsByDay: O(n)
// - Инициализация первого окна: O(vacationLength) = O(p)
// - Скользящее окно: O(p - vacationLength) = O(p)
// - Итого: O(n + p)
//
// Пространственная сложность: O(p)
// - Массив meetingsByDay: O(p)
// - Остальные переменные: O(1)
// - Итого: O(p)
//
// Почему скользящее окно эффективнее?
// - Наивный подход: для каждого окна пересчитывать сумму → O(p * vacationLength)
// - Скользящее окно: обновляем сумму за O(1) → O(p)
// - Ускорение в vacationLength раз!
//
// Для vacationLength = 14 и periodLength = 31:
// - Наивный: 18 * 14 = 252 операции
// - Скользящее окно: 31 операций
// - Ускорение в ~14 раз
