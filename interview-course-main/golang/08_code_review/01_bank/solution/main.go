package main

// СПИСОК НАЙДЕННЫХ ПРОБЛЕМ И ИСПРАВЛЕНИЙ:
//
// 1. SQL инъекция (строка 119 оригинала)
//    Проблема: Параметры подставляются через fmt.Sprintf
//    Решение: Использовать параметризованные запросы db.Exec с плейсхолдерами
//
// 2. Хардкод credentials (строки 94-100)
//    Проблема: Пароли и хосты в коде, попадут в git
//    Решение: Читать из переменных окружения или конфиг-файла
//
// 3. Новое подключение к БД на каждый запрос (строка 105)
//    Проблема: sql.Open при каждом вызове updateCurrency
//    Решение: Создать connection pool один раз при старте
//
// 4. Неправильный defer в цикле (строка 68)
//    Проблема: defer resp.Body.Close() накапливается в цикле
//    Решение: Вызывать Close() явно или вынести в функцию
//
// 5. Игнорирование ошибок (строки 58, 69)
//    Проблема: Ошибки игнорируются через _
//    Решение: Проверять и обрабатывать все ошибки
//
// 6. Panic вместо обработки ошибок (строки 66, 80, 85, 107, 114)
//    Проблема: panic() при ошибке, приложение падает
//    Решение: Логировать ошибку и продолжать или возвращать ошибку наверх
//
// 7. Нет конкурентности (строка 57)
//    Проблема: Запросы выполняются последовательно в цикле for
//    Решение: Использовать горутины и sync.WaitGroup
//
// 8. Нет таймаутов для HTTP запросов
//    Проблема: Запрос может висеть вечно
//    Решение: Добавить timeout в http.Client
//
// 9. Нет graceful shutdown
//    Проблема: Нет обработки сигналов ОС
//    Решение: Использовать signal.Notify и context для отмены
//
// 10. Нет контекста для cancel'а операций
//     Проблема: Нельзя отменить долгие операции
//     Решение: Использовать context.Context
//
// 11. Нет повторных попыток при ошибках
//     Проблема: Одна ошибка - вся операция провалена
//     Решение: Добавить retry logic с exponential backoff
//
// 12. fmt.Println для логов (строка 117)
//     Проблема: Неструктурированное логирование
//     Решение: Использовать logger (zap, zerolog, slog)
//
// 13. Нет валидации данных
//     Проблема: Парсим float без проверки формата ответа
//     Решение: Валидировать структуру ответа
//
// 14. Неправильный параметр в updateCurrency (строка 82)
//     Проблема: Передаем url.curFrom дважды вместо curFrom и curTo
//     Решение: Исправить на url.curFrom, url.curTo
//
// 15. Нет rate limiting
//     Проблема: Можем заDDoS'ить банковские API
//     Решение: Добавить rate limiter
//
// 16. Нет метрик и мониторинга
//     Проблема: Не видим что происходит в проде
//     Решение: Добавить prometheus метрики
//
// 17. Нет транзакций для БД операций
//     Проблема: Нет транзакций для БД операций
//     Решение: Использовать db.BeginTx
//
// 18. db.Ping() на каждый запрос (строка 112)
//     Проблема: Лишняя нагрузка, достаточно проверить при старте
//     Решение: Проверять только при инициализации
//
// ПРАВИЛЬНАЯ РЕАЛИЗАЦИЯ приведена ниже в виде концепции.
// В реальном проекте нужно добавить больше обработки ошибок,
// тесты, конфигурацию и т.д.

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	_ "github.com/lib/pq"
)

// Config хранит конфигурацию приложения
type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
}

// LoadConfig загружает конфигурацию из переменных окружения
func LoadConfig() *Config {
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     5432,
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "currency"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// InitDB инициализирует connection pool к БД
func InitDB(cfg *Config) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	// Проверяем подключение один раз при старте
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	// Настраиваем connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

// UpdateCurrency сохраняет курс валюты в БД (правильная версия)
func UpdateCurrency(ctx context.Context, db *sql.DB, bank, from, to string, value float64) error {
	// Используем параметризованный запрос для защиты от SQL injection
	query := `INSERT INTO currency_rates (bank, "from", "to", value, created_at)
	          VALUES ($1, $2, $3, $4, NOW())`

	_, err := db.ExecContext(ctx, query, bank, from, to, value)
	if err != nil {
		return fmt.Errorf("failed to insert currency: %w", err)
	}

	return nil
}

// BankRate представляет информацию о курсе валюты для банка
type BankRate struct {
	BankName  string
	CurFrom   string
	CurTo     string
	URL       string
	AuthToken string // Опционально для банков, требующих авторизацию
}

// fetchRate получает курс валюты от банка
func fetchRate(ctx context.Context, client *http.Client, rate BankRate) (float64, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rate.URL, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	// Добавляем авторизацию если требуется
	if rate.AuthToken != "" {
		req.Header.Add("Authorization", rate.AuthToken)
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Читаем body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read body: %w", err)
	}

	strBody := string(body)

	// Некоторые банки возвращают числа с запятой вместо точки
	if rate.BankName == "Bank 1" {
		strBody = strings.ReplaceAll(strBody, ",", ".")
	}

	// Парсим float
	value, err := strconv.ParseFloat(strings.TrimSpace(strBody), 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse value '%s': %w", strBody, err)
	}

	return value, nil
}

// processBank обрабатывает получение и сохранение курса для одного банка
func processBank(ctx context.Context, db *sql.DB, client *http.Client, rate BankRate) error {
	log.Printf("Fetching rate from %s...", rate.BankName)

	value, err := fetchRate(ctx, client, rate)
	if err != nil {
		return fmt.Errorf("%s: %w", rate.BankName, err)
	}

	log.Printf("%s: %s/%s = %.4f", rate.BankName, rate.CurFrom, rate.CurTo, value)

	err = UpdateCurrency(ctx, db, rate.BankName, rate.CurFrom, rate.CurTo, value)
	if err != nil {
		return fmt.Errorf("%s: failed to save: %w", rate.BankName, err)
	}

	return nil
}

// updateAllRates обновляет курсы валют для всех банков параллельно
func updateAllRates(ctx context.Context, db *sql.DB, rates []BankRate) error {
	// HTTP клиент с таймаутами
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			IdleConnTimeout:     30 * time.Second,
			TLSHandshakeTimeout: 5 * time.Second,
		},
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(rates))

	// Запускаем горутины для каждого банка
	for _, rate := range rates {
		wg.Add(1)
		go func(r BankRate) {
			defer wg.Done()

			if err := processBank(ctx, db, client, r); err != nil {
				errChan <- err
			}
		}(rate)
	}

	// Ждем завершения всех горутин
	wg.Wait()
	close(errChan)

	// Собираем ошибки
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
		log.Printf("Error: %v", err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to update %d rates", len(errors))
	}

	return nil
}

func run() error {
	// Загружаем конфигурацию
	cfg := LoadConfig()

	// Инициализируем БД
	db, err := InitDB(cfg)
	if err != nil {
		return fmt.Errorf("failed to init db: %w", err)
	}
	defer db.Close()

	log.Println("Database connected successfully")

	// Настраиваем graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Обрабатываем сигналы для graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Received signal: %v, shutting down gracefully...", sig)
		cancel()
	}()

	// Список банков для опроса
	rates := []BankRate{
		{
			BankName: "Bank 1",
			CurFrom:  "RUB",
			CurTo:    "USD",
			URL:      "http://bank.example.com/rates/rub-usd",
		},
		{
			BankName:  "Bank 2",
			CurFrom:   "RUB",
			CurTo:     "USD",
			URL:       "http://bank2.example.com/rates?currFrom=RUR&currTo=USD",
			AuthToken: "Bearer XXXXXXX",
		},
	}

	// Обновляем курсы
	log.Println("Starting currency rates update...")
	err = updateAllRates(ctx, db, rates)
	if err != nil {
		return fmt.Errorf("failed to update rates: %w", err)
	}

	log.Println("All rates updated successfully")
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
