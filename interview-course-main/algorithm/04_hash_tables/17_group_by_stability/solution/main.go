package main

import (
	"fmt"
	"strconv"
)

// ============================================================================
// Решение: Группировка по ключу через map
// ============================================================================

func groupByStability(stats []ServerStat) map[string][]int {
	result := make(map[string][]int)

	for _, stat := range stats {
		// Преобразуем float64 в string для использования как ключ
		key := strconv.FormatFloat(stat.Stability, 'f', -1, 64)

		// Добавляем server ID в список для данной стабильности
		result[key] = append(result[key], stat.Server)
	}

	return result
}

// ============================================================================
// Демонстрация и тесты
// ============================================================================

func main() {
	fmt.Println("=== Решение: Группировка по стабильности ===")
	testGroupByStability()
}

func testGroupByStability() {
	stats := []ServerStat{
		{Server: 1, Stability: 99},
		{Server: 2, Stability: 97},
		{Server: 3, Stability: 34},
		{Server: 4, Stability: 97},
		{Server: 5, Stability: 97.1},
	}

	result := groupByStability(stats)
	fmt.Printf("Результат: %v\n", result)

	// Проверка ожидаемого результата
	expected := map[string][]int{
		"34":   {3},
		"97":   {2, 4},
		"99":   {1},
		"97.1": {5},
	}

	status := "✅"
	for key, expectedList := range expected {
		if actualList, ok := result[key]; !ok || !slicesEqual(actualList, expectedList) {
			status = "❌"
			break
		}
	}
	fmt.Printf("%s Тест пройден\n", status)
}

func slicesEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// ============================================================================
// Объяснение решения
// ============================================================================

// РЕШЕНИЕ: Группировка через map
//
// Идея:
// - Используем map, где ключ - стабильность (строка), значение - список ID серверов
// - Проходим по всем статистикам и группируем серверы по значению стабильности
//
// Алгоритм:
// 1. Создаем пустой map[string][]int
// 2. Для каждой статистики:
//    - Преобразуем Stability (float64) в string (ключ)
//    - Добавляем Server ID в список для этого ключа
// 3. Возвращаем map
//
// Пример работы:
// stats = [
//   {Server: 1, Stability: 99},
//   {Server: 2, Stability: 97},
//   {Server: 3, Stability: 34},
//   {Server: 4, Stability: 97},
//   {Server: 5, Stability: 97.1}
// ]
//
// Итерация 1: Server=1, Stability=99
//   key = "99"
//   result["99"] = [1]
//
// Итерация 2: Server=2, Stability=97
//   key = "97"
//   result["97"] = [2]
//
// Итерация 3: Server=3, Stability=34
//   key = "34"
//   result["34"] = [3]
//
// Итерация 4: Server=4, Stability=97
//   key = "97"
//   result["97"] = [2, 4]  // добавляем к существующему списку
//
// Итерация 5: Server=5, Stability=97.1
//   key = "97.1"
//   result["97.1"] = [5]
//
// Результат: {"34": [3], "97": [2, 4], "99": [1], "97.1": [5]}
//
// Временная сложность: O(n)
// - n - количество элементов в stats
// - Проходим по массиву один раз: O(n)
// - Операции с map (append, доступ по ключу): O(1) amortized
// - Итого: O(n)
//
// Пространственная сложность: O(n)
// - В худшем случае все серверы имеют разные стабильности
// - Map будет содержать n ключей, каждый с одним элементом
// - Итого: O(n) памяти
//
// Почему strconv.FormatFloat?
// - float64 нельзя использовать напрямую как ключ map
// - Нужно преобразовать в string для использования как ключ
// - 'f' - формат без экспоненты
// - -1 - использовать минимальное количество знаков после запятой
// - 64 - размер float64
//
// Альтернатива: можно использовать map[float64][]int
// - Но float64 как ключ может привести к проблемам с точностью
// - Например: 97.0 и 97.0000000001 будут разными ключами
// - String ключи более надежны для группировки
