package main

import "fmt"

// Задача: Удаление дубликатов из отсортированного списка
//
// Решение: Однопроходный алгоритм с текущим указателем
//
// Идея:
// 1. Проходим по списку с указателем current
// 2. Сравниваем current.Val с current.Next.Val
// 3. Если значения равны - пропускаем дубликат (current.Next = current.Next.Next)
// 4. Если значения разные - переходим к следующему узлу (current = current.Next)
//
// Пример: [1,1,2,3,3]
// Шаг 1: current=1, next=1 (равны) → пропускаем → [1,2,3,3]
// Шаг 2: current=1, next=2 (разные) → переходим → current=2
// Шаг 3: current=2, next=3 (разные) → переходим → current=3
// Шаг 4: current=3, next=3 (равны) → пропускаем → [1,2,3]
// Шаг 5: current=3, next=nil → завершаем
//
// Важно:
// - Список уже отсортирован, поэтому дубликаты всегда идут подряд
// - Не нужна дополнительная память для хеш-таблицы
// - Изменяем список in-place
//
// Сложность:
// - Время: O(n) - один проход по списку
// - Память: O(1) - используем только указатель current

type ListNode struct {
	Val  int
	Next *ListNode
}

func deleteDuplicates(head *ListNode) *ListNode {
	// Пустой список или один элемент
	if head == nil {
		return head
	}

	current := head

	// Проходим по всему списку
	for current != nil && current.Next != nil {
		// Если текущее значение равно следующему
		if current.Val == current.Next.Val {
			// Пропускаем дубликат
			current.Next = current.Next.Next
		} else {
			// Переходим к следующему узлу
			current = current.Next
		}
	}

	return head
}

func main() {
	// Тест 1: [1,1,2] → [1,2]
	head := &ListNode{1, &ListNode{1, &ListNode{2, nil}}}
	fmt.Print("Before: ")
	printList(head)
	result := deleteDuplicates(head)
	fmt.Print("After: ")
	printList(result) // Ожидается: 1 -> 2
	fmt.Println()

	// Тест 2: [1,1,2,3,3] → [1,2,3]
	head = &ListNode{1, &ListNode{1, &ListNode{2, &ListNode{3, &ListNode{3, nil}}}}}
	fmt.Print("Before: ")
	printList(head)
	result = deleteDuplicates(head)
	fmt.Print("After: ")
	printList(result) // Ожидается: 1 -> 2 -> 3
	fmt.Println()

	// Тест 3: [1,1,1,1] → [1]
	head = &ListNode{1, &ListNode{1, &ListNode{1, &ListNode{1, nil}}}}
	fmt.Print("Before: ")
	printList(head)
	result = deleteDuplicates(head)
	fmt.Print("After: ")
	printList(result) // Ожидается: 1
	fmt.Println()

	// Тест 4: [1,2,3] (нет дубликатов) → [1,2,3]
	head = &ListNode{1, &ListNode{2, &ListNode{3, nil}}}
	fmt.Print("Before: ")
	printList(head)
	result = deleteDuplicates(head)
	fmt.Print("After: ")
	printList(result) // Ожидается: 1 -> 2 -> 3
}

func printList(head *ListNode) {
	if head == nil {
		fmt.Println("Empty list")
		return
	}
	for head != nil {
		fmt.Print(head.Val)
		if head.Next != nil {
			fmt.Print(" -> ")
		}
		head = head.Next
	}
	fmt.Println()
}
