package main

import "fmt"

// MergeToMap добавляет уникальные значения из newValues в слайс по ключу key.
func MergeToMap(m map[string][]string, key string, newValues []string) {
	// Получаем существующие значения для ключа.
	existing := m[key]

	// Создаем map для быстрой проверки существующих значений.
	existingSet := make(map[string]bool)
	for _, v := range existing {
		existingSet[v] = true
	}

	// Добавляем только те значения, которых еще нет.
	for _, v := range newValues {
		if !existingSet[v] {
			existing = append(existing, v)
			existingSet[v] = true // Обновляем set для предотвращения дубликатов в newValues
		}
	}

	// Обновляем мапу.
	m[key] = existing
}

func main() {
	m := map[string][]string{
		"group1": {"apple", "banana"},
		"group2": {"carrot"},
	}

	newValues := []string{"banana", "cherry"}
	key := "group1"

	fmt.Println("До MergeToMap:", m)
	MergeToMap(m, key, newValues)
	fmt.Println("После MergeToMap:", m)
}

// Ответ:
// До MergeToMap: map[group1:[apple banana] group2:[carrot]]
// После MergeToMap: map[group1:[apple banana cherry] group2:[carrot]]
//
// Объяснение:
// 1. Создаем map для O(1) проверки существования значений.
// 2. Проходим по newValues и добавляем только уникальные значения.
// 3. Обновляем мапу с новым слайсом.
//
// Оптимизация: использование map для проверки уникальности дает O(n) сложность
// вместо O(n²) при проверке через слайс.
