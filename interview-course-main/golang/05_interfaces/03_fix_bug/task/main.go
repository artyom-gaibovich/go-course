// Задача: Исправьте код так, чтобы:
//  1. Реализации интерфейсов работали корректно.
//  2. Программа завершалась успешно и без паники.
//  3. Лишние вызовы методов были устранены.
//  4. Вы смогли объяснить, где были ошибки и почему они возникли.

package main

import (
	"errors"
	"fmt"
)

// PaymentProcessor - интерфейс для обработки платежей
type PaymentProcessor interface {
	Process(amount float64) error
	Verify(amount float64) bool
}

// CreditCardProcessor - реализация интерфейса для кредитной карты
type CreditCardProcessor struct {
	limit float64
}

// Process обрабатывает платеж
func (c *CreditCardProcessor) Process(amount float64) error {
	if amount > c.limit {
		return errors.New("credit card limit exceeded")
	}
	fmt.Printf("Process payment of $%.2f using CreditCard\n", amount)
	return nil
}

// Verify проверяет возможность обработки платежа
func (c CreditCardProcessor) Verify(amount float64) bool {
	return amount <= c.limit
}

// PayPalProcessor - реализация интерфейса для PayPal
type PayPalProcessor struct {
	balance float64
}

// Process обрабатывает платеж
func (p *PayPalProcessor) Process(amount float64) error {
	if amount > p.balance {
		return errors.New("not enough balance in PayPal")
	}
	fmt.Printf("Processed payment of $%.2f using PayPal\n", amount)
	return nil
}

// Verify проверяет возможность обработки платежа
func (p *PayPalProcessor) Verify(amount float64) bool {
	return amount <= p.balance
}

// ExecutePayment вызывает методы Process и Verify
func ExecutePayment(processor PaymentProcessor, amount float64) {
	if processor.Verify(amount) {
		err := processor.Process(amount)
		if err != nil {
			fmt.Println("Error:", err)
		}
	} else {
		fmt.Println("Verification failed for amount:", amount)
	}
}

func main() {
	creditCard := CreditCardProcessor{limit: 100.0}
	payPal := PayPalProcessor{balance: 200.0}

	ExecutePayment(creditCard, 50.0)
	ExecutePayment(&creditCard, 50.0)
	ExecutePayment(&payPal, 150.0)
	ExecutePayment(payPal, 150.0)
}
