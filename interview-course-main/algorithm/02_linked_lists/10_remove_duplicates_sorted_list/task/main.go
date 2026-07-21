package main

// Задача: Удаление дубликатов из отсортированного списка
//
// Дан отсортированный по возрастанию связный список.
// Необходимо удалить все дубликаты так, чтобы каждый элемент встречался только один раз.
//
// Требования:
// - Сложность по времени O(n)
// - Сложность по памяти O(1)
//
// Примеры:
// Input: head = [1,1,2]
// Output: [1,2]
//
// Input: head = [1,1,2,3,3]
// Output: [1,2,3]
//
// Input: head = [1,1,1,1]
// Output: [1]

type ListNode struct {
	Val  int
	Next *ListNode
}

func deleteDuplicates(head *ListNode) *ListNode {
	// TODO: реализуйте функцию
	return nil
}

func main() {
	// Тест 1: [1,1,2] → [1,2]
	head := &ListNode{1, &ListNode{1, &ListNode{2, nil}}}
	result := deleteDuplicates(head)
	printList(result) // Ожидается: 1 -> 2

	// Тест 2: [1,1,2,3,3] → [1,2,3]
	head = &ListNode{1, &ListNode{1, &ListNode{2, &ListNode{3, &ListNode{3, nil}}}}}
	result = deleteDuplicates(head)
	printList(result) // Ожидается: 1 -> 2 -> 3

	// Тест 3: [1,1,1,1] → [1]
	head = &ListNode{1, &ListNode{1, &ListNode{1, &ListNode{1, nil}}}}
	result = deleteDuplicates(head)
	printList(result) // Ожидается: 1
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
