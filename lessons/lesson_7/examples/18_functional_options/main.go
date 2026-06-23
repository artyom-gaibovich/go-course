package main

import "fmt"

type User struct {
	Name    string
	Surname string
	Email   string
	Phone   string
}

type Option func(*User)

func WithEmail(email string) Option {
	return func(u *User) {
		u.Email = email
	}
}

func WithPhone(phone string) Option {
	return func(u *User) {
		u.Phone = phone
	}
}

func NewUser(name, surname string, options ...Option) User {
	user := User{Name: name, Surname: surname}
	for _, option := range options {
		option(&user)
	}
	return user
}

func main() {
	u1 := NewUser("Ivan", "Ivanov", WithEmail("ivan@example.com"))
	u2 := NewUser("Petr", "Petrov", WithEmail("petr@example.com"), WithPhone("+700000000"))

	fmt.Printf("%+v\n%+v\n", u1, u2)
}
