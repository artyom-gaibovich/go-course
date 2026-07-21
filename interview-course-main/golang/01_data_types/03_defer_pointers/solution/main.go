package main

import "fmt"

type Account struct {
	Balance int
}

func main() {
	initialBalance := 1000
	account := &Account{Balance: initialBalance}

	// defer вычисляет аргументы в момент вызова defer.
	// account.Balance передается как значение (int), поэтому копируется СЕЙЧАС (1000).
	defer printBalance("Изначальный баланс", account.Balance)

	// account.Balance снова копируется в момент вызова defer (все еще 1000).
	defer printBalance("Текущий баланс", account.Balance)

	// account передается как указатель, но сам указатель копируется в момент вызова defer.
	// Указатель указывает на ту же структуру, что и account.
	defer printAccountBalance("Указатель на баланс", account)

	// Изменяем баланс через указатель.
	account.Balance += 500 // account.Balance = 1500

	// Изменяем баланс через функцию, которая принимает указатель.
	updateBalance(account, 200) // account.Balance = 1300

	// ВАЖНО: здесь мы переназначаем переменную account на новый объект.
	// Но defer уже захватил старый указатель, поэтому это не влияет на отложенные вызовы.
	account = &Account{Balance: 300}
}

func updateBalance(acc *Account, amount int) {
	acc.Balance -= amount
}

func printBalance(label string, amount int) {
	fmt.Printf("%s: %d\n", label, amount)
}

func printAccountBalance(label string, acc *Account) {
	fmt.Printf("%s: %d\n", label, acc.Balance)
}

// Ответ:
// Указатель на баланс: 1300
// Текущий баланс: 1000
// Изначальный баланс: 1000
//
// Объяснение:
// 1. defer выполняются в обратном порядке (LIFO).
// 2. printBalance получает значение account.Balance, которое копируется в момент вызова defer (1000).
// 3. printAccountBalance получает указатель на account, который указывает на ту же структуру.
//    Когда мы изменяем account.Balance, это влияет на значение по указателю.
// 4. Переназначение account = &Account{Balance: 300} не влияет на defer,
//    так как defer уже захватил старый указатель.
