package main

// Задача: Find All Numbers Disappeared in an Array
//
// Дан массив nums из n целых чисел, где nums[i] находится в диапазоне [1, n].
// Верните массив всех целых чисел в диапазоне [1, n], которые не появляются в nums.
//
// Входные данные:
// - nums: []int - массив целых чисел, где каждый элемент в диапазоне [1, n]
//
// Выходные данные:
// - []int - массив всех чисел из диапазона [1, n], которых нет в nums
//
// Требования:
// - n == nums.length
// - 1 <= n <= 10^5
// - 1 <= nums[i] <= n
//
// Дополнительное требование: Решить без дополнительной памяти и за O(n) времени.
// Возвращаемый список не считается дополнительной памятью.
//
// Примеры:
// Input: nums = [4,3,2,7,8,2,3,1]
// Output: [5,6]
//
// Input: nums = [1,1]
// Output: [2]

func findDisappearedNumbers(nums []int) []int {
	// TODO: реализуйте функцию
	return nil
}

func main() {
	// Тест 1
	nums1 := []int{4, 3, 2, 7, 8, 2, 3, 1}
	result1 := findDisappearedNumbers(nums1)
	// Ожидается: [5,6]
	println("Test 1:", result1)

	// Тест 2
	nums2 := []int{1, 1}
	result2 := findDisappearedNumbers(nums2)
	// Ожидается: [2]
	println("Test 2:", result2)
}
