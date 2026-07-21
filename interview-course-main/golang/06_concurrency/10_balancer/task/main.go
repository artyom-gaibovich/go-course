// Задача:
// Реализовать client-side балансировщик нагрузки для микросервисов.
//
// Контекст:
// - Есть микросервис с интерфейсом Backend
// - Для каждого микросервиса запущено несколько десятков экземпляров
// - Каждый экземпляр доступен по своему адресу
// - Экземпляры ненадежны: могут падать, быть недоступными или перегруженными
//
// Требования:
// - Реализовать тип Balancer, который также реализует Backend
// - Использовать алгоритм Round Robin для балансировки
// - Учитывать конкурентный доступ (метод может вызываться из горутин)
// - Балансировщик должен равномерно распределять нагрузку

package main

import "context"

type Request interface{}

type Response interface{}

type Backend interface {
	Invoke(ctx context.Context, req Request) (Response, error)
}

type BackendImpl struct {
	addr string
}

var _ Backend = &BackendImpl{}

// NewBackend создает backend для конкретного экземпляра по адресу.
func NewBackend(addr string) *BackendImpl {
	return &BackendImpl{addr: addr}
}

// Invoke отправляет запрос на конкретный backend.
func (b *BackendImpl) Invoke(ctx context.Context, req Request) (Response, error) {
	// Заглушка: в реальности здесь HTTP/gRPC вызов
	return "response from " + b.addr, nil
}

type Balancer struct {
	// TODO: добавьте необходимые поля
}

var _ Backend = &Balancer{}

// NewBalancer создает балансировщик для списка адресов backend'ов.
func NewBalancer(addrs []string) *Balancer {
	panic("implement me")
}

// Invoke распределяет запросы между backend'ами по Round Robin.
func (bal *Balancer) Invoke(ctx context.Context, req Request) (Response, error) {
	panic("implement me")
}
