package main

// Задача: Удаление N-го элемента с конца односвязного списка
//
// Дан односвязный список и число n. Нужно удалить n-й элемент с конца списка
// и вернуть голову списка.
//
// Требования:
// - Решить задачу за один проход по списку
//
// Примеры:
// Input: head = [1,2,3,4,5], n = 2
// Output: [1,2,3,5]
// Объяснение: Удаляем 4 (второй элемент с конца)
//
// Input: head = [1], n = 1
// Output: []
// Объяснение: Удаляем единственный элемент
//
// Input: head = [1,2], n = 1
// Output: [1]
// Объяснение: Удаляем последний элемент

type ListNode struct {
	Val  int
	Next *ListNode
}

func removeNthFromEnd(head *ListNode, n int) *ListNode {
	// TODO: реализуйте функцию
	return nil
}

func main() {
	// Тест 1: [1,2,3,4,5], n=2 → [1,2,3,5]
	head := &ListNode{1, &ListNode{2, &ListNode{3, &ListNode{4, &ListNode{5, nil}}}}}
	result := removeNthFromEnd(head, 2)
	printList(result) // Ожидается: 1 -> 2 -> 3 -> 5

	// Тест 2: [1], n=1 → []
	head = &ListNode{1, nil}
	result = removeNthFromEnd(head, 1)
	printList(result) // Ожидается: пустой список

	// Тест 3: [1,2], n=1 → [1]
	head = &ListNode{1, &ListNode{2, nil}}
	result = removeNthFromEnd(head, 1)
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
