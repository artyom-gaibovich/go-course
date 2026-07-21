package main

// Задача: Разворот связного списка
//
// Дан односвязный список. Необходимо развернуть его.
//
// Требования:
// - Реализовать итеративное решение
// - Сложность по времени O(n)
// - Сложность по памяти O(1)
//
// Примеры:
// Input: head = [1,2,3,4,5]
// Output: [5,4,3,2,1]
//
// Input: head = [1,2]
// Output: [2,1]
//
// Input: head = []
// Output: []

type ListNode struct {
	Val  int
	Next *ListNode
}

func reverseList(head *ListNode) *ListNode {
	// TODO: реализуйте функцию
	return nil
}

func main() {
	// Тест 1: [1,2,3,4,5] → [5,4,3,2,1]
	head := &ListNode{1, &ListNode{2, &ListNode{3, &ListNode{4, &ListNode{5, nil}}}}}
	result := reverseList(head)
	printList(result) // Ожидается: 5 -> 4 -> 3 -> 2 -> 1

	// Тест 2: [1,2] → [2,1]
	head = &ListNode{1, &ListNode{2, nil}}
	result = reverseList(head)
	printList(result) // Ожидается: 2 -> 1

	// Тест 3: [] → []
	result = reverseList(nil)
	printList(result) // Ожидается: пустой список
}

func printList(head *ListNode) {
	if head == nil {
		println("Empty list")
		return
	}
	for head != nil {
		print(head.Val)
		if head.Next != nil {
			print(" -> ")
		}
		head = head.Next
	}
	println()
}
