/*
Домашнее задание 1: Собственный буферизированный канал.

Реализуйте обобщённый тип MyChan[T any] поверх sync.Mutex + sync.Cond и
циклического буфера. Он должен повторять поведение встроенного буферизированного
канала:

  - NewMyChan[T](size int) *MyChan[T] — создать канал с буфером на size элементов.
  - Send(v T)        — положить элемент; блокируется, пока буфер полон;
                       паникует, если канал закрыт.
  - Recv() (T, bool) — забрать элемент; блокируется, пока буфер пуст и канал открыт;
                       после закрытия сначала отдаёт остаток буфера (ok=true),
                       затем zero value, ok=false.
  - Close()          — закрыть канал; повторное закрытие должно паниковать.

Критерии:
  [ ] корректные блокировки на границах полного/пустого буфера;
  [ ] паника при Send в закрытый канал и при повторном Close;
  [ ] чтение из закрытого канала отдаёт остаток буфера, затем ok=false;
  [ ] go run -race не показывает гонок;
  [ ] нет утечек горутин.

Подсказка: используйте два sync.Cond (notFull, notEmpty) над одним мьютексом,
поля sendx/recvx/qcount для циклического буфера и флаг closed.
*/

package main

type MyChan[T any] struct {
	// TODO: mu, notFull, notEmpty, buf, sendx, recvx, qcount, closed
}

func NewMyChan[T any](size int) *MyChan[T] {
	// TODO
	return nil
}

func (c *MyChan[T]) Send(v T) {
	// TODO
}

func (c *MyChan[T]) Recv() (T, bool) {
	// TODO
	var zero T
	return zero, false
}

func (c *MyChan[T]) Close() {
	// TODO
}

func main() {
	// TODO: продемонстрировать работу MyChan (send/recv, close, drain после close)
}
