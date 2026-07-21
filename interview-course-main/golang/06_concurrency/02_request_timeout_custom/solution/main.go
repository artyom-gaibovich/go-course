package main

import (
	"context"
	"fmt"
	"time"
)

// Эта функция лезет по сети в старый монолит и может тупить.
func getDiscount(ctx context.Context) (float64, error) {
	// Симулируем сетевой запрос
	select {
	case <-time.After(2 * time.Second):
		return 12.0, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

func main() {
	// Создаем контекст с таймаутом 1 секунда
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	discount, err := getDiscount(ctx)
	if err != nil {
		fmt.Printf("Ошибка получения скидки: %v. Используем скидку по умолчанию: 0\n", err)
		discount = 0
	}

	fmt.Printf("Ваша скидка: %v\n", discount)
}

// Альтернативное решение с каналом:
// func getDiscountWithTimeout(timeout time.Duration) (float64, error) {
//     resultChan := make(chan float64, 1)
//     errorChan := make(chan error, 1)
//
//     go func() {
//         discount := getDiscount()
//         resultChan <- discount
//     }()
//
//     select {
//     case discount := <-resultChan:
//         return discount, nil
//     case <-time.After(timeout):
//         return 0, fmt.Errorf("timeout after %v", timeout)
//     }
// }

// Объяснение:
// 1. Добавляем context.Context для контроля таймаута.
// 2. Используем select для ожидания либо результата, либо таймаута.
// 3. Обрабатываем ошибку и предоставляем значение по умолчанию.
// 4. Это предотвращает блокировку всей программы на долгий сетевой запрос.
