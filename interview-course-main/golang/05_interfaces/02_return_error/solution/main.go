package main

func main() {
	println(handle())
}

// В Go error - это интерфейс: type error interface { Error() string }
// Можно вернуть nil (нет ошибки) или создать свою реализацию error.
func handle() error {
	// Вариант 1: вернуть nil (нет ошибки)
	// return nil

	// Вариант 2: создать свою реализацию error без импорта пакетов
	return &myError{msg: "произошла ошибка"}
}

// myError - простая реализация интерфейса error.
type myError struct {
	msg string
}

func (e *myError) Error() string {
	return e.msg
}

// Объяснение:
// 1. error - это интерфейс с одним методом Error() string.
// 2. Можно вернуть nil (нет ошибки).
// 3. Можно создать свою структуру, реализующую интерфейс error.
// 4. Не нужно импортировать пакеты errors или fmt для базовой реализации.
