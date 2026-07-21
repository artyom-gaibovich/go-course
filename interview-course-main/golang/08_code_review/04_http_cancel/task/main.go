// Задача: Сделать ревью кода и добавить отмену запросов при ошибке
//
// ТЗ: программа отправляет HTTP GET запросы к нескольким URL параллельно.
// Если хотя бы один запрос завершается с ошибкой, нужно отменить все остальные.
//
// Требования:
// - Запросы должны выполняться параллельно
// - При первой ошибке все остальные запросы должны быть отменены
//
// Этот код НАМЕРЕННО содержит ошибки для учебных целей!
// Не запускайте в production!

package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func main() {
	urls := []string{
		"https://google.com",
		"https://yandex.ru",
		"https://github.com",
		"https://stackoverflow.com",
	}

	for _, url := range urls {
		go func(url string) {
			err := fetch(context.Background(), url)
			if err != nil {
				fmt.Printf("Error fetching %s: %v\n", url, err)
				return
			}
			fmt.Printf("Success: %s\n", url)
		}(url)
	}

	fmt.Println("All requests launched!")
	time.Sleep(400 * time.Millisecond)
	fmt.Println("Done")
}

func fetch(ctx context.Context, url string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	return nil
}
