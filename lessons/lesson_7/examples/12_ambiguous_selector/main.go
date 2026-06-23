package main

import "fmt"

type Derived1 struct {
	Values int
}

func (d Derived1) Print() {
	fmt.Println("derived1:", d.Values)
}

type Derived2 struct {
	Values int
}

func (d Derived2) Print() {
	fmt.Println("derived2:", d.Values)
}

type Data struct {
	Derived1
	Derived2
}

func main() {
	data := Data{
		Derived1: Derived1{Values: 1},
		Derived2: Derived2{Values: 2},
	}

	fmt.Println(data.Derived1.Values)
	fmt.Println(data.Derived2.Values)

	data.Derived1.Print()
	data.Derived2.Print()
}
