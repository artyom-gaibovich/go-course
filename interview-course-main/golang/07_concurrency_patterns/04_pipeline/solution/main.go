package main

import "fmt"

type Transaction struct {
	ID       int
	Amount   float64
	Currency string
}

// ReadStage читает транзакции из источника и отправляет в выходной канал.
func ReadStage(transactions []Transaction) <-chan Transaction {
	out := make(chan Transaction)
	go func() {
		defer close(out)
		fmt.Println("[ReadStage] Начало чтения транзакций")
		for i, t := range transactions {
			fmt.Printf("[ReadStage] Читаю транзакцию %d: ID=%d, Amount=%.2f %s\n", i+1, t.ID, t.Amount, t.Currency)
			out <- t
		}
		fmt.Println("[ReadStage] Завершено чтение всех транзакций")
	}()
	return out
}

// FilterStage фильтрует транзакции с отрицательными суммами.
// Принимает канал от предыдущей стадии, возвращает канал для следующей.
func FilterStage(in <-chan Transaction) <-chan Transaction {
	out := make(chan Transaction)
	go func() {
		defer close(out)
		fmt.Println("[FilterStage] Начало фильтрации транзакций")
		for t := range in {
			if t.Amount > 0 {
				fmt.Printf("[FilterStage] Транзакция ID=%d прошла фильтр (Amount=%.2f > 0)\n", t.ID, t.Amount)
				out <- t
			} else {
				fmt.Printf("[FilterStage] Транзакция ID=%d отфильтрована (Amount=%.2f <= 0)\n", t.ID, t.Amount)
			}
		}
		fmt.Println("[FilterStage] Завершена фильтрация")
	}()
	return out
}

// ConvertStage конвертирует валюту в доллары.
// Принимает канал от предыдущей стадии, возвращает канал для следующей.
func ConvertStage(in <-chan Transaction) <-chan Transaction {
	out := make(chan Transaction)
	go func() {
		defer close(out)
		fmt.Println("[ConvertStage] Начало конвертации валют")
		rates := map[string]float64{
			"EUR": 1.1,
			"RUB": 0.011,
			"USD": 1.0,
		}
		for t := range in {
			oldAmount := t.Amount
			oldCurrency := t.Currency
			rate := rates[t.Currency]
			t.Amount = t.Amount * rate
			t.Currency = "USD"
			fmt.Printf("[ConvertStage] Конвертирована транзакция ID=%d: %.2f %s -> %.2f %s (rate=%.3f)\n",
				t.ID, oldAmount, oldCurrency, t.Amount, t.Currency, rate)
			out <- t
		}
		fmt.Println("[ConvertStage] Завершена конвертация")
	}()
	return out
}

// SaveStage сохраняет результаты.
// Принимает канал от предыдущей стадии, собирает результаты в слайс.
func SaveStage(in <-chan Transaction) []Transaction {
	var results []Transaction
	fmt.Println("[SaveStage] Начало сохранения результатов")
	count := 0
	for t := range in {
		count++
		fmt.Printf("[SaveStage] Сохранена транзакция %d: ID=%d, Amount=%.2f %s\n", count, t.ID, t.Amount, t.Currency)
		results = append(results, t)
	}
	fmt.Printf("[SaveStage] Завершено сохранение. Всего сохранено: %d транзакций\n", count)
	return results
}

func main() {
	transactions := []Transaction{
		{ID: 1, Amount: 100, Currency: "EUR"},
		{ID: 2, Amount: -50, Currency: "USD"},
		{ID: 3, Amount: 200, Currency: "RUB"},
		{ID: 4, Amount: 150, Currency: "USD"},
	}

	// Строим pipeline: каждая стадия принимает канал от предыдущей и возвращает канал для следующей.
	stage1 := ReadStage(transactions) // Чтение транзакций
	stage2 := FilterStage(stage1)     // Фильтрация (принимает канал от stage1)
	stage3 := ConvertStage(stage2)    // Конвертация (принимает канал от stage2)
	results := SaveStage(stage3)      // Сохранение (принимает канал от stage3)

	fmt.Println("Обработанные транзакции:")
	for _, t := range results {
		fmt.Printf("ID: %d, Amount: %.2f %s\n", t.ID, t.Amount, t.Currency)
	}
}

// Объяснение:
// 1. Pipeline паттерн обрабатывает данные через несколько этапов.
// 2. Каждая стадия принимает канал от предыдущей стадии и возвращает канал для следующей.
// 3. Стадии работают параллельно: данные передаются через каналы по мере обработки.
// 4. Каждая стадия запускается в отдельной горутине и обрабатывает данные потоком.
// 5. Pipeline позволяет легко добавлять/удалять стадии и комбинировать их.
//
// Преимущества такого подхода:
// - Каждая стадия независима и может быть переиспользована.
// - Легко тестировать каждую стадию отдельно.
// - Можно легко добавлять новые стадии в цепочку.
// - Параллельная обработка ускоряет общую производительность.
