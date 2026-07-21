// Задача: Что выведет программа и почему?

package main

import "fmt"

type User struct {
	Name string
}

func main() {
	user := User{Name: "Гриша"}
	fmt.Println("Имя до обновления:", user.Name)

	updateUser(user)
	fmt.Println("Имя после обновления:", user.Name)
}

func updateUser(u User) {
	u.Name = "Таненбаум"
	fmt.Println("Имя внутри функции [updateUser]:", u.Name)

	resetUser(&u)
	fmt.Println("Имя после вызова [resetUser] внутри функции [updateUser]:", u.Name)
}

func resetUser(u *User) {
	u = &User{Name: "Безымянный"}
	fmt.Println("Имя внутри функции [resetUser]:", u.Name)
}
