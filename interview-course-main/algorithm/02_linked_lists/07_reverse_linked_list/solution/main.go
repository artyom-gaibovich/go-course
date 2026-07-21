package main

import "fmt"

// Задача: Разворот связного списка
//
// Решение: Итеративный подход с тремя указателями
//
// Идея:
// Используем три указателя: prev, current, next
// 1. prev = nil (новый конец списка)
// 2. current = head (текущий узел)
// 3. Для каждого узла:
//    - Сохраняем next = current.Next
//    - Меняем направление: current.Next = prev
//    - Сдвигаем указатели: prev = current, current = next
//
// Пример: 1 -> 2 -> 3 -> nil
// Шаг 1: nil <- 1    2 -> 3 -> nil   (prev=1, current=2)
// Шаг 2: nil <- 1 <- 2    3 -> nil   (prev=2, current=3)
// Шаг 3: nil <- 1 <- 2 <- 3          (prev=3, current=nil)
// Результат: 3 -> 2 -> 1 -> nil
//
// Альтернативное решение: рекурсивный подход
// - Рекурсивно разворачиваем остаток списка
// - Меняем направление текущего узла
// - Возвращаем новую голову
//
// Сложность:
// - Время: O(n) - один проход по списку
// - Память: O(1) для итеративного, O(n) для рекурсивного (стек вызовов)

type ListNode struct {
	Val  int
	Next *ListNode
}

// Итеративное решение
func reverseList(head *ListNode) *ListNode {
	var prev *ListNode
	current := head

	for current != nil {
		// Сохраняем следующий узел
		next := current.Next

		// Разворачиваем указатель
		current.Next = prev

		// Сдвигаем указатели
		prev = current
		current = next
	}

	// prev указывает на новую голову (бывший последний элемент)
	return prev
}

// Рекурсивное решение (альтернативный подход)
func reverseListRecursive(head *ListNode) *ListNode {
	// Базовый случай: пустой список или последний элемент
	if head == nil || head.Next == nil {
		return head
	}

	// Рекурсивно разворачиваем остаток списка
	newHead := reverseListRecursive(head.Next)

	// Меняем направление текущего узла
	head.Next.Next = head
	head.Next = nil

	return newHead
}

func main() {
	// Тест 1: [1,2,3,4,5] → [5,4,3,2,1]
	head := &ListNode{1, &ListNode{2, &ListNode{3, &ListNode{4, &ListNode{5, nil}}}}}
	fmt.Print("Before: ")
	printList(head)
	result := reverseList(head)
	fmt.Print("After (iterative): ")
	printList(result)
	fmt.Println()

	// Тест 2: [1,2] → [2,1]
	head = &ListNode{1, &ListNode{2, nil}}
	fmt.Print("Before: ")
	printList(head)
	result = reverseList(head)
	fmt.Print("After (iterative): ")
	printList(result)
	fmt.Println()

	// Тест 3: рекурсивное решение
	head = &ListNode{1, &ListNode{2, &ListNode{3, nil}}}
	fmt.Print("Before: ")
	printList(head)
	result = reverseListRecursive(head)
	fmt.Print("After (recursive): ")
	printList(result)
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
