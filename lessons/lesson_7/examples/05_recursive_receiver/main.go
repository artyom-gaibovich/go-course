package main

import "fmt"

type Account struct {
	balance int
}

type Client struct {
	account *Account
}

func (c Client) Deposit(v int) {
	c.account.balance += v
}

func main() {
	client := Client{account: &Account{}}

	client.Deposit(100)
	client.Deposit(50)

	fmt.Println("balance:", client.account.balance)
}
