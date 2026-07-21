// Задача: Что выведет программа и почему?

package main

import "fmt"

type person struct {
	age int
}

func main() {
	people := make([]person, 2)

	p1 := &people[1]
	fmt.Printf("%p", p1)

	p1.age++

	people = append(people, person{}, person{}, person{})
	fmt.Println(cap(people))

	p1.age++

	fmt.Println(people[1].age)
	fmt.Println(p1.age)
}
