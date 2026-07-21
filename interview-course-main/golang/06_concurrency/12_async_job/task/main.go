// Задача:
// MultiProcess выполняет input во всех переданных JobFunc параллельно.
// Как только одна из JobFunc успешно завершается, результат возвращается.
// Если все JobFunc завершаются с ошибкой, возвращается последняя ошибка.

package main

import (
	"context"
)

type Result struct{}

// JobFunc - асинхронная задача, которая принимает контекст и параметр,
// выполняет обработку и возвращает результат или ошибку.
type JobFunc func(ctx context.Context, input string) (Result, error)

func MultiProcess(ctx context.Context, input string, jobs []JobFunc) (Result, error) {
	// TODO: реализуйте параллельное выполнение с возвратом первого успеха
	panic("implement me")
}
