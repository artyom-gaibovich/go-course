/*
Задание 3: Система флагов прав доступа пользователя

Реализуй систему прав через битовые флаги (аналог chmod).

Флаги:
  PermRead  = 1 << 0  // чтение
  PermWrite = 1 << 1  // запись
  PermExec  = 1 << 2  // выполнение
  PermAdmin = 1 << 3  // администратор

Реализуй функции:
  Grant(perms, flag int) int        — выдать право
  Revoke(perms, flag int) int       — отозвать право
  Has(perms, flag int) bool         — проверить наличие права
  String(perms int) string          — строка вида "rw-a" или "r---"

Формат String: позиции r, w, x, a (в этом порядке). Если флаг не установлен — дефис.
Пример: Grant(0, PermRead|PermExec) → String → "r-x-"

Подсказки:
  Grant  → OR  (|)
  Revoke → AND NOT (&^)
  Has    → AND и сравнение с флагом (не с 0, чтобы корректно работать с составными масками)
*/

package main

import "fmt"

const (
	PermRead  = 1 << 0
	PermWrite = 1 << 1
	PermExec  = 1 << 2
	PermAdmin = 1 << 3
)

func Grant(perms, flag int) int {
	// TODO
	return perms
}

func Revoke(perms, flag int) int {
	// TODO
	return perms
}

func Has(perms, flag int) bool {
	// TODO
	return false
}

func String(perms int) string {
	// TODO: вернуть строку вида "rw-a"
	return "----"
}

func main() {
	perms := 0
	perms = Grant(perms, PermRead|PermWrite)
	fmt.Printf("After Grant(Read|Write): %s\n", String(perms))

	perms = Grant(perms, PermAdmin)
	fmt.Printf("After Grant(Admin):      %s\n", String(perms))

	perms = Revoke(perms, PermWrite)
	fmt.Printf("After Revoke(Write):     %s\n", String(perms))

	fmt.Printf("Has(Read):  %v\n", Has(perms, PermRead))
	fmt.Printf("Has(Write): %v\n", Has(perms, PermWrite))
	fmt.Printf("Has(Admin): %v\n", Has(perms, PermAdmin))
}
