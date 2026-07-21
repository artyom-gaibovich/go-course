// Задача: Что выведет программа и почему?

package main

import "fmt"

func main() {
	nums := []int{1, 2, 3}

	addNum(nums[0:2])
	fmt.Println(nums)

	addNums(nums[0:2])
	fmt.Println(nums)
}

func addNum(nums []int) {
	nums = append(nums, 4)
}

func addNums(nums []int) {
	nums = append(nums, 5, 6)
}
