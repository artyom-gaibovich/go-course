package main

import "fmt"

const (
	OperatorPlus     = "+"
	OperatorMinus    = "-"
	OperatorMultiply = "*"
)

func calculator(operator string, a, b int) {
	if operator == OperatorPlus {
		fmt.Println(a + b)
	}
	if operator == OperatorMinus {
		fmt.Println(a - b)
	}
	if operator == OperatorMultiply {
		fmt.Println(a * b)
	}
}

func calculator2(operator string, a, b int) {
	if operator == OperatorPlus {
		fmt.Println(a + b*10)
	}
	if operator == OperatorMinus {
		fmt.Println(a - b*10)
	}
	if operator == OperatorMultiply {
		fmt.Println(a * b * 10)
	}
}

func main() {
	a, b := 10, 20

	calculator("+", a, b)

}
