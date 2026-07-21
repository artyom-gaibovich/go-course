package main

import "fmt"

// Задача: Поиск цикла в связном списке
//
// Решение: Алгоритм Floyd's Cycle Detection (Черепаха и Заяц)
//
// Идея:
// Используем два указателя:
// 1. slow (медленный) - движется на 1 шаг за итерацию
// 2. fast (быстрый) - движется на 2 шага за итерацию
//
// Если есть цикл:
// - Быстрый указатель в конце концов догонит медленный
// - Они встретятся в какой-то точке внутри цикла
//
// Если нет цикла:
// - Быстрый указатель достигнет конца списка (nil)
//
// Почему это работает:
// - Если есть цикл, оба указателя окажутся внутри цикла
// - Разница в скоростях = 1 шаг за итерацию
// - Быстрый указатель "догоняет" медленный на 1 шаг каждую итерацию
// - Встреча гарантирована
//
// Сложность:
// - Время: O(n) - в худшем случае медленный указатель проходит весь список
// - Память: O(1) - используем только два указателя

type ListNode struct {
	Val  int
	Next *ListNode
}

func hasCycle(head *ListNode) bool {
	// Пустой список или один элемент без цикла
	if head == nil || head.Next == nil {
		return false
	}

	// Инициализируем указатели
	slow := head
	fast := head

	// Пока быстрый указатель может двигаться
	for fast != nil && fast.Next != nil {
		// Медленный делает 1 шаг
		slow = slow.Next

		// Быстрый делает 2 шага
		fast = fast.Next.Next

		// Если указатели встретились - есть цикл
		if slow == fast {
			return true
		}
	}

	// Быстрый указатель достиг конца - нет цикла
	return false
}

func main() {
	// Тест 1: [3,2,0,-4] с циклом на позиции 1
	node4 := &ListNode{Val: -4}
	node3 := &ListNode{Val: 0, Next: node4}
	node2 := &ListNode{Val: 2, Next: node3}
	head := &ListNode{Val: 3, Next: node2}
	node4.Next = node2                                 // создаем цикл: -4 -> 2
	fmt.Println("Test 1 (цикл есть):", hasCycle(head)) // Ожидается: true

	// Тест 2: [1,2] с циклом на позиции 0
	node22 := &ListNode{Val: 2}
	head2 := &ListNode{Val: 1, Next: node22}
	node22.Next = head2                                 // создаем цикл: 2 -> 1
	fmt.Println("Test 2 (цикл есть):", hasCycle(head2)) // Ожидается: true

	// Тест 3: [1] без цикла
	head3 := &ListNode{Val: 1}
	fmt.Println("Test 3 (нет цикла):", hasCycle(head3)) // Ожидается: false

	// Тест 4: [1,2,3,4,5] без цикла
	head4 := &ListNode{Val: 1, Next: &ListNode{Val: 2, Next: &ListNode{Val: 3, Next: &ListNode{Val: 4, Next: &ListNode{Val: 5}}}}}
	fmt.Println("Test 4 (нет цикла):", hasCycle(head4)) // Ожидается: false
}
