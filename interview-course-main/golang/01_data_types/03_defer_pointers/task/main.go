// Задача: Что выведет программа и почему?

package main

import "fmt"

type Account struct {
	Balance int
}

func main() {
	initialBalance := 1000
	account := &Account{Balance: initialBalance} // Balance = 1000

	defer printBalance("Изначальный баланс", account.Balance) // 1000
	defer printAccountBalance("Указатель на баланс", account) // oldStruct = 1300

	account.Balance += 500           // 1500
	updateBalance(account, 200)      // 1300
	account = &Account{Balance: 300} // oldStruct = 1300, newStruct = 300
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
