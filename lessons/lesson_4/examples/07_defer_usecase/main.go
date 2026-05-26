package main

import (
	"fmt"
	"os"
)

// Главный use case defer — гарантированная очистка ресурсов.
// Открыл — сразу отложил закрытие, чтобы не забыть при любом исходе.

func readFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close() // выполнится при любом выходе из функции

	buf := make([]byte, 64)
	n, err := f.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

// defer + замыкание: можно изменить именованное возвращаемое значение
func safeDiv(a, b int) (result int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("перехвачена паника: %v", r)
		}
	}()
	result = a / b // паника при b == 0
	return
}

func main() {
	content, err := readFile("/etc/hostname")
	if err != nil {
		fmt.Println("Ошибка чтения:", err)
	} else {
		fmt.Println("Хост:", content)
	}

	res, err := safeDiv(10, 2)
	fmt.Printf("10/2 = %d, err = %v\n", res, err)

	res, err = safeDiv(10, 0)
	fmt.Printf("10/0 = %d, err = %v\n", res, err)
}
