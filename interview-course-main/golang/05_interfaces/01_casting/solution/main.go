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
	var userReader Reader = &User{}
	// Type assertion: проверяем, реализует ли userReader интерфейс Writer.
	// Это возможно, так как &User{} реализует оба интерфейса (Reader и Writer).
	userWriter := userReader.(Writer)
	userWriter.Get()
	_ = userWriter
}

// Ответ: программа компилируется и выполняется без ошибок.
//
// Объяснение:
// 1. &User{} реализует оба интерфейса: Reader (методы Get, List) и Writer (методы Create, Delete).
// 2. Type assertion userReader.(Writer) успешен, так как userReader содержит *User,
//    который реализует интерфейс Writer.
// 3. userWriter имеет тип Writer, но указывает на тот же объект *User.
// 4. Вызов userWriter.Get() корректен, так как Writer не запрещает вызов других методов.
//    Интерфейс определяет минимальный набор методов, но объект может иметь больше методов.
//
// Важно: интерфейс в Go определяет набор методов, которые объект ДОЛЖЕН реализовать,
// но не ограничивает объект только этими методами.
