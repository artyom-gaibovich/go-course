package main

import "fmt"

type Person struct {
	Name string
}

func (p Person) Intro() {
	fmt.Println("I am", p.Name)
}

type Woman struct {
	Person
}

func (w Woman) Intro() {
	fmt.Println("Mrs.", w.Person.Name)
}

func main() {
	woman := Woman{Person: Person{Name: "Ekaterina"}}

	woman.Intro()
	woman.Person.Intro()
}
