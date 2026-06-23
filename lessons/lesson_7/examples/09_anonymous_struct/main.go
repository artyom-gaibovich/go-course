package main

import "fmt"

func main() {
	type localRequest struct {
		ID   int
		Name string
	}

	req := localRequest{ID: 1, Name: "local"}
	fmt.Println(req)

	book := struct {
		Title  string
		Author struct {
			Name string
		}
	}{
		Title: "Go",
		Author: struct {
			Name string
		}{
			Name: "Ken",
		},
	}

	fmt.Println(book.Title, book.Author.Name)
}
