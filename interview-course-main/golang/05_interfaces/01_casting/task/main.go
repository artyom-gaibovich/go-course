// Задача: Что выведет программа и почему?

package main

type User struct{}

func (u *User) Create() {}
func (u *User) Get()    {}
func (u *User) List()   {}
func (u *User) Delete() {}

type Reader interface {
	Get()
	List()
}

type Writer interface {
	Create()
	Delete()
}

func main() {
	// var userReader Reader = &User{}
	// userWriter := userReader.(Writer)
	// userWriter.Get()
	// _ = userWriter
}
