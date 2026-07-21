package main

// Задача: Слияние двух отсортированных связных списков
//
// Даны два отсортированных по возрастанию связных списка.
// Необходимо слить их в один отсортированный список.
//
// Требования:
// - Сложность по времени O(n + m)
// - Сложность по памяти O(1) - не создавать новые узлы
//
// Примеры:
// Input: list1 = [1,2,4], list2 = [1,3,4]
// Output: [1,1,2,3,4,4]
//
// Input: list1 = [], list2 = []
// Output: []
//
// Input: list1 = [], list2 = [0]
// Output: [0]

type ListNode struct {
	Val  int
	Next *ListNode
}

func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
	// TODO: реализуйте функцию
	return nil
}

func main() {
	// Тест 1: [1,2,4] и [1,3,4] → [1,1,2,3,4,4]
	list1 := &ListNode{1, &ListNode{2, &ListNode{4, nil}}}
	list2 := &ListNode{1, &ListNode{3, &ListNode{4, nil}}}
	result := mergeTwoLists(list1, list2)
	printList(result) // Ожидается: 1 -> 1 -> 2 -> 3 -> 4 -> 4

	// Тест 2: [] и [] → []
	result = mergeTwoLists(nil, nil)
	printList(result) // Ожидается: пустой список

	// Тест 3: [] и [0] → [0]
	list2 = &ListNode{0, nil}
	result = mergeTwoLists(nil, list2)
	printList(result) // Ожидается: 0
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
