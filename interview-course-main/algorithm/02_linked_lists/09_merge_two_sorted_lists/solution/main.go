package main

import "fmt"

// Задача: Слияние двух отсортированных связных списков
//
// Решение: Итеративный подход с dummy node
//
// Идея:
// 1. Создаем dummy node (фиктивный узел) для упрощения логики
// 2. Используем указатель current для построения результата
// 3. Сравниваем текущие узлы list1 и list2
// 4. Добавляем меньший узел к результату
// 5. После завершения одного списка, добавляем остаток другого
//
// Пример: list1 = [1,2,4], list2 = [1,3,4]
// Шаг 1: dummy -> 1 (из list1), list1=[2,4], list2=[1,3,4]
// Шаг 2: dummy -> 1 -> 1 (из list2), list1=[2,4], list2=[3,4]
// Шаг 3: dummy -> 1 -> 1 -> 2 (из list1), list1=[4], list2=[3,4]
// Шаг 4: dummy -> 1 -> 1 -> 2 -> 3 (из list2), list1=[4], list2=[4]
// Шаг 5: dummy -> 1 -> 1 -> 2 -> 3 -> 4 (из list1), list1=[], list2=[4]
// Шаг 6: dummy -> 1 -> 1 -> 2 -> 3 -> 4 -> 4 (из list2)
//
// Альтернативное решение: рекурсивный подход
// - Базовый случай: если один из списков пуст, возвращаем другой
// - Сравниваем головы списков
// - Рекурсивно сливаем остаток
//
// Сложность:
// - Время: O(n + m) - проходим по обоим спискам один раз
// - Память: O(1) для итеративного, O(n + m) для рекурсивного (стек вызовов)

type ListNode struct {
	Val  int
	Next *ListNode
}

// Итеративное решение
func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
	// Создаем dummy node для упрощения логики
	dummy := &ListNode{}
	current := dummy

	// Пока оба списка не пусты
	for list1 != nil && list2 != nil {
		if list1.Val <= list2.Val {
			current.Next = list1
			list1 = list1.Next
		} else {
			current.Next = list2
			list2 = list2.Next
		}
		current = current.Next
	}

	// Добавляем остаток (один из списков уже закончился)
	if list1 != nil {
		current.Next = list1
	} else {
		current.Next = list2
	}

	return dummy.Next
}

// Рекурсивное решение (альтернативный подход)
func mergeTwoListsRecursive(list1 *ListNode, list2 *ListNode) *ListNode {
	// Базовые случаи
	if list1 == nil {
		return list2
	}
	if list2 == nil {
		return list1
	}

	// Сравниваем головы и рекурсивно сливаем остаток
	if list1.Val <= list2.Val {
		list1.Next = mergeTwoListsRecursive(list1.Next, list2)
		return list1
	} else {
		list2.Next = mergeTwoListsRecursive(list1, list2.Next)
		return list2
	}
}

func main() {
	// Тест 1: [1,2,4] и [1,3,4] → [1,1,2,3,4,4]
	list1 := &ListNode{1, &ListNode{2, &ListNode{4, nil}}}
	list2 := &ListNode{1, &ListNode{3, &ListNode{4, nil}}}
	fmt.Print("Test 1 (iterative): ")
	result := mergeTwoLists(list1, list2)
	printList(result) // Ожидается: 1 -> 1 -> 2 -> 3 -> 4 -> 4

	// Тест 2: [] и [] → []
	fmt.Print("Test 2 (iterative): ")
	result = mergeTwoLists(nil, nil)
	printList(result) // Ожидается: пустой список

	// Тест 3: рекурсивное решение
	list1 = &ListNode{1, &ListNode{3, &ListNode{5, nil}}}
	list2 = &ListNode{2, &ListNode{4, &ListNode{6, nil}}}
	fmt.Print("Test 3 (recursive): ")
	result = mergeTwoListsRecursive(list1, list2)
	printList(result) // Ожидается: 1 -> 2 -> 3 -> 4 -> 5 -> 6
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
