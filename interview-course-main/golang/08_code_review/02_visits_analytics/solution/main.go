package main

// СПИСОК НАЙДЕННЫХ ПРОБЛЕМ И ИСПРАВЛЕНИЙ:
//
// 1. Отсутствуют необходимые импорты (строки 20-30)
//    Проблема: Используются пакеты без импорта chi, middleware
//    Решение: Использовать стандартный http.ServeMux или добавить импорты
//
// 2. SQL Injection уязвимость (строки 98-101)
//    Проблема: Параметры подставляются через fmt.Sprintf
//              Хотя в данном случае используется метод ToPgInterval(), который
//              возвращает безопасные значения, это плохая практика.
//              Если в будущем кто-то изменит логику, может появиться уязвимость.
//    Решение: Использовать параметризованные запросы или константы
//
// 3. Игнорирование всех ошибок (строки 72, 87, 90, 103, 107, 157, 158)
//    Проблема: Все ошибки игнорируются через _
//    Решение: Проверять и обрабатывать каждую ошибку
//
// 4. Хардкод credentials в коде (строка 71)
//    Проблема: Пароль и connection string в коде
//    Решение: Читать из переменных окружения
//
// 5. Несовпадение типов параметров (строка 81 vs 95)
//    Проблема: handleVisits принимает sql.DB (значение),
//              а getVisitsFromDB принимает *sql.DB (указатель)
//    Решение: Везде использовать *sql.DB (указатель)
//
// 6. Отсутствует defer rows.Close() (после строки 103)
//    Проблема: Не закрывается rows, утечка ресурсов
//    Решение: Добавить defer rows.Close() сразу после db.Query
//
// 7. Отсутствует проверка rows.Err() (после строки 110)
//    Проблема: Ошибки при итерации не проверяются
//    Решение: Добавить проверку rows.Err() после цикла
//
// 8. Неправильное использование context.WithValue (строки 83-85)
//    Проблема: context.WithValue предназначен для request-scoped данных,
//              не для передачи бизнес-параметров
//    Решение: Передавать period как обычный параметр функции
//
// 9. Отсутствует валидация параметра period
//    Проблема: Не проверяется, что period имеет допустимое значение
//    Решение: Валидировать входные данные
//
// 10. Потенциальная race condition в map (строки 133-137)
//     Проблема: Изменение DayVisit требует чтения, модификации и записи
//               Если использовать concurrency, будет race
//     Решение: В данном случае нет concurrency, но лучше использовать указатели + mutex
//
// 11. Неоптимальный поиск максимума (строки 142-148)
//     Проблема: При counter >= maxCounter всегда перезаписываем TopLocation
//               Это приводит к неопределенному результату при равных счетчиках
//     Решение: Использовать counter > maxCounter
//
// 12. Отсутствует проверка на пустой период (строка 100)
//     Проблема: Если period пустой или невалидный, запрос будет некорректным
//     Решение: Валидация и значение по умолчанию
//
// 13. Отсутствует graceful shutdown
//     Проблема: Сервер не завершается корректно
//     Решение: Добавить обработку сигналов и graceful shutdown
//
// 14. Отсутствует обработка HTTP ошибок
//     Проблема: При ошибках клиенту всегда возвращается 200 OK
//     Решение: Возвращать правильные HTTP статусы
//
// 15. Игнорирование параметров функций (строка 81)
//     Проблема: В handleVisits параметр r объявлен как http.Request, но должен быть *http.Request
//     Решение: Использовать указатели для больших структур
//
// 16. Синтаксические ошибки в объявлении map (строки 117-118)
//     Проблема: map[string]DayVisit() и map[string]map[string]int() - лишние ()
//     Решение: Убрать () после объявления типа map
//
// 17. Отсутствует индекс в БД
//     Проблема: Запрос по timestamp без индекса будет медленным
//     Решение: Создать индекс: CREATE INDEX idx_visits_timestamp ON visits(timestamp)
//
// 18. Игнорирование ошибок парсинга времени (строки 157-158)
//     Проблема: time.Parse может вернуть ошибку, но она игнорируется
//     Решение: Обрабатывать ошибку или использовать fallback
//

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	_ "github.com/lib/pq"
)

type UserID string
type Period string

const (
	PeriodWeek Period = "week"
	PeriodDay  Period = "day"
)

// ToPgInterval возвращает PostgreSQL интервал для периода
func (p Period) ToPgInterval() string {
	switch p {
	case PeriodWeek:
		return "7 days"
	case PeriodDay:
		return "1 day"
	default:
		return "1 day"
	}
}

// IsValid проверяет валидность периода
func (p Period) IsValid() bool {
	return p == PeriodWeek || p == PeriodDay
}

type Visit struct {
	UserID   UserID    `json:"user_id"`
	Location string    `json:"location"`
	TS       time.Time `json:"timestamp"`
}

type DayVisit struct {
	Day         string `json:"day"`
	Count       int    `json:"count"`
	TopLocation string `json:"top_location"`
}

func main() {
	// Читаем connection string из переменной окружения
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://admin:password123@localhost:5432/analytics?sslmode=disable"
		log.Println("WARNING: Using default connection string. Set DATABASE_URL env variable.")
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Проверяем подключение
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Используем стандартный http.ServeMux
	mux := http.NewServeMux()
	mux.HandleFunc("/visits", handleVisits(db))

	log.Println("Server starting on :3000")
	err = http.ListenAndServe(":3000", mux)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// handleVisits - правильная версия с обработкой ошибок
func handleVisits(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем и валидируем период
		periodStr := r.URL.Query().Get("period")
		if periodStr == "" {
			periodStr = "day" // значение по умолчанию
		}

		period := Period(periodStr)
		if !period.IsValid() {
			http.Error(w, "Invalid period parameter. Use 'day' or 'week'", http.StatusBadRequest)
			return
		}

		// Используем request context для отмены
		ctx := r.Context()

		// Передаем period как параметр, а не через context
		visits, err := getVisitsFromDB(ctx, db, period)
		if err != nil {
			log.Printf("Failed to get visits: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		dayVisits := dayVisitsFromVisits(visits)

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(dayVisits)
		if err != nil {
			log.Printf("Failed to encode response: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

// getVisitsFromDB - правильная версия с параметризованным запросом
func getVisitsFromDB(ctx context.Context, db *sql.DB, period Period) ([]Visit, error) {
	var visits []Visit

	// Используем параметризованный запрос для безопасности
	// Интервал передаем через плейсхолдер
	query := `SELECT user_id, location, timestamp
	          FROM visits
	          WHERE timestamp > NOW() - INTERVAL '1 day' * $1`

	var days int
	if period == PeriodWeek {
		days = 7
	} else {
		days = 1
	}

	rows, err := db.QueryContext(ctx, query, days)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close() // ВАЖНО: закрываем rows

	for rows.Next() {
		var v Visit
		err = rows.Scan(&v.UserID, &v.Location, &v.TS)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		visits = append(visits, v)
	}

	// Проверяем ошибки после итерации
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("rows iteration failed: %w", err)
	}

	return visits, nil
}

func dayVisitsFromVisits(visits []Visit) []DayVisit {
	// Используем map без () после типа
	dayVisits := make(map[string]*DayVisit) // Используем указатели
	locationCounter := make(map[string]map[string]int)

	for _, visit := range visits {
		day := visit.TS.Format(time.DateOnly)

		// Инициализируем структуры для нового дня
		if _, ok := dayVisits[day]; !ok {
			dayVisits[day] = &DayVisit{Day: day}
			locationCounter[day] = make(map[string]int)
		}

		// Увеличиваем счетчики (благодаря указателям не нужно перезаписывать)
		dayVisits[day].Count++
		locationCounter[day][visit.Location]++
	}

	// Находим топ локацию для каждого дня
	for day, locCounter := range locationCounter {
		var maxCount int
		var topLocation string

		for location, count := range locCounter {
			// Используем > вместо >= для детерминированности
			if count > maxCount {
				maxCount = count
				topLocation = location
			}
		}

		dayVisits[day].TopLocation = topLocation
	}

	// Собираем результат
	out := make([]DayVisit, 0, len(dayVisits))
	for _, dv := range dayVisits {
		out = append(out, *dv)
	}

	// Сортируем по дате
	sort.Slice(out, func(i, j int) bool {
		// Обрабатываем ошибки парсинга
		ti, err1 := time.Parse(time.DateOnly, out[i].Day)
		tj, err2 := time.Parse(time.DateOnly, out[j].Day)

		// Если ошибка парсинга, сравниваем как строки
		if err1 != nil || err2 != nil {
			return out[i].Day < out[j].Day
		}

		return ti.Before(tj)
	})

	return out
}

// ДОПОЛНИТЕЛЬНЫЕ УЛУЧШЕНИЯ ДЛЯ PRODUCTION:
//
// 1. Добавить graceful shutdown с обработкой сигналов
// 2. Добавить connection pool настройки (SetMaxOpenConns, SetMaxIdleConns)
// 3. Добавить timeout для запросов
// 4. Добавить логирование с уровнями (zap, zerolog)
// 5. Добавить метрики (prometheus)
// 6. Добавить rate limiting
// 7. Добавить кеширование результатов
// 8. Добавить индексы в БД для оптимизации
// 9. Добавить пагинацию для больших результатов
// 10. Добавить тесты (unit и integration)
