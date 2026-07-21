// Задача: Что выведет программа и почему?

package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	str := "ddЯЙ異"
	fmt.Println("Длина через len:", len(str))
	fmt.Println("Длина через RuneCountInString:", utf8.RuneCountInString(str))
}
