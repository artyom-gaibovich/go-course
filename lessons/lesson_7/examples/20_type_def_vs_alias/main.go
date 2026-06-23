package main

import "fmt"

type (
	AliasInt = int
	MyInt    int
)

func (m MyInt) Double() MyInt {
	return m * 2
}

func main() {
	var a AliasInt = 10
	var plain int = a
	fmt.Println("alias assigns without cast:", plain)

	var m MyInt = 21
	var n int = int(m)
	fmt.Println("definition needs cast:", n)

	fmt.Println("method on defined type:", m.Double())
}
