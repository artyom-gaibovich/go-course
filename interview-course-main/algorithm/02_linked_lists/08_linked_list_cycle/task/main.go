package main

// Задача: Поиск цикла в связном списке
//
// Дан односвязный список. Необходимо определить, есть ли в нём цикл.
//
// Требования:
// - Реализовать алгоритм Floyd's Cycle Detection (медленный и быстрый указатели)
// - Сложность по времени O(n)
// - Сложность по памяти O(1)
//
// Примеры:
// Input: head = [3,2,0,-4], pos = 1 (цикл начинается с узла со значением 2)
// Output: true
//
// Input: head = [1,2], pos = 0 (цикл начинается с узла со значением 1)
// Output: true
//
// Input: head = [1], pos = -1 (нет цикла)
// Output: false

type ListNode struct {
	Val  int
	Next *ListNode
}

func hasCycle(head *ListNode) bool {
	// TODO: реализуйте функцию
	return false
}

func main() {
	// Тест 1: [3,2,0,-4] с циклом на позиции 1
	node4 := &ListNode{Val: -4}
	node3 := &ListNode{Val: 0, Next: node4}
	node2 := &ListNode{Val: 2, Next: node3}
	head := &ListNode{Val: 3, Next: node2}
	node4.Next = node2                 // создаем цикл
	println("Test 1:", hasCycle(head)) // Ожидается: true

	// Тест 2: [1,2] с циклом на позиции 0
	node22 := &ListNode{Val: 2}
	head2 := &ListNode{Val: 1, Next: node22}
	node22.Next = head2                 // создаем цикл
	println("Test 2:", hasCycle(head2)) // Ожидается: true

	// Тест 3: [1] без цикла
	head3 := &ListNode{Val: 1}
	println("Test 3:", hasCycle(head3)) // Ожидается: false
}
