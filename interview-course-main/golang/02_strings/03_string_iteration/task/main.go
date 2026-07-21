// Задача: Что выведет программа и почему?

package main

import "fmt"

func main() {
	str := "G😀o"
	for i := 0; i < len(str); i++ {
		fmt.Printf("%c ", str[i])
	}
}
