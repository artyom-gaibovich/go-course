package main

import "fmt"

type Node struct {
	value int
	next  *Node
}

func (n *Node) Describe() string {
	if n == nil {
		return "nil node"
	}
	return fmt.Sprintf("node(%d)", n.value)
}

func main() {
	var empty *Node
	fmt.Println(empty.Describe())

	filled := &Node{value: 5}
	fmt.Println(filled.Describe())
}
