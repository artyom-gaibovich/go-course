package main

import "fmt"

type Entity struct {
	Field1 int
	Field2 string
	Field3 int
}

type OODStorage struct {
	entities []Entity
}

type DODStorage struct {
	field1 []int
	field2 []string
	field3 []int
}

func main() {
	const n = 5

	ood := OODStorage{entities: make([]Entity, n)}
	for i := range ood.entities {
		ood.entities[i] = Entity{Field1: i, Field2: "x", Field3: i * 2}
	}

	dod := DODStorage{
		field1: make([]int, n),
		field2: make([]string, n),
		field3: make([]int, n),
	}
	for i := 0; i < n; i++ {
		dod.field1[i] = i
		dod.field2[i] = "x"
		dod.field3[i] = i * 2
	}

	oodSum := 0
	for i := range ood.entities {
		oodSum += ood.entities[i].Field1
	}

	dodSum := 0
	for i := range dod.field1 {
		dodSum += dod.field1[i]
	}

	fmt.Println("OOD sum:", oodSum, "DOD sum:", dodSum)
}
