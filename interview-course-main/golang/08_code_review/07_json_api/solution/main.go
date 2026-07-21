package main

// СПИСОК НАЙДЕННЫХ ПРОБЛЕМ И ИСПРАВЛЕНИЙ:
//
// 1. Отсутствует ограничение размера request body (строка 46)
//    Проблема: io.ReadAll(r.Body) читает весь body в память
//              Атакующий может отправить 10GB JSON → out of memory
//    Решение: http.MaxBytesReader для ограничения размера
//
// 2. Игнорируются все ошибки (строки 46, 50, 58, 69, 75)
//    Проблема: json.Unmarshal, ReadAll возвращают errors, которые игнорируются
//              Невалидный JSON обработается как пустая структура
//    Решение: Проверять все ошибки и возвращать HTTP 400
//
// 3. Отсутствует Content-Type проверка (строка 41)
//    Проблема: API принимает любой Content-Type
//              Можно отправить text/plain вместо application/json
//    Решение: Проверять Content-Type header
//
// 4. Нет валидации входных данных (строка 49-58)
//    Проблема: Принимается любой JSON
//              Name может быть пустым, Age = -100, Email без @
//    Решение: Валидировать поля перед сохранением
//
// 5. Race condition при доступе к map (строки 55, 78)
//    Проблема: Конкурентный доступ к map без синхронизации
//    Решение: Использовать sync.RWMutex или sync.Map
//
// 6. Race condition для nextID (строка 54)
//    Проблема: nextID++ не атомарная операция
//              Два запроса могут получить одинаковый ID
//    Решение: Использовать atomic.AddInt64
//
// 7. Неправильные HTTP статусы
//    Проблема: Всегда возвращается 200 OK
//              Создание должно возвращать 201 Created
//              Ошибки должны возвращать 4xx/5xx
//              Пустой результат должен возвращать 404
//    Решение: Использовать правильные статусы
//
// 8. Отсутствует обработка несуществующего пользователя (строка 78)
//    Проблема: GET /user/999 вернет пустую структуру со статусом 200
//    Решение: Проверять существование и возвращать 404
//
// 9. Response не содержит Content-Type
//     Проблема: Content-Type не указан, клиент не знает формат ответа
//     Решение: Устанавливать Content-Type: application/json
//
// 10. Отсутствует обработка defer body.Close()
//     Проблема: Request body не закрывается
//     Решение: Добавить defer r.Body.Close()
//
// 11. JSON декодер создается напрямую для записи
//     Проблема: Если Encode() вернет ошибку, она игнорируется
//     Решение: Проверять ошибку или использовать json.Marshal
//
// 12. Небезопасный парсинг ID из URL (строка 75)
//     Проблема: fmt.Sscanf может не распарсить, id останется 0
//     Решение: Использовать роутер или проверять результат
//
// 13. Нет rate limiting
//     Проблема: Атакующий может залить сервер запросами
//     Решение: Добавить rate limiter
//
// 14. Нет логирования запросов
//     Проблема: Не видно что происходит
//     Решение: Middleware для логирования
//
// 15. Глобальные переменные users и nextID (строки 25-26)
//     Проблема: Сложно тестировать, нет инкапсуляции
//     Решение: Создать UserService с методами
//
// 16. Panic при конкурентных запросах к map
//     Проблема: concurrent map read and map write
//     Решение: Синхронизация доступа

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

// User представляет пользователя
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// Validate валидирует пользователя
func (u *User) Validate() error {
	if strings.TrimSpace(u.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if len(u.Name) > 100 {
		return fmt.Errorf("name is too long (max 100 characters)")
	}
	if !strings.Contains(u.Email, "@") {
		return fmt.Errorf("invalid email format")
	}
	if u.Age < 0 || u.Age > 150 {
		return fmt.Errorf("age must be between 0 and 150")
	}
	return nil
}

// ErrorResponse представляет ответ с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
}

// UserService управляет пользователями
type UserService struct {
	mu     sync.RWMutex
	users  map[int]User
	nextID int64
}

// NewUserService создает новый сервис
func NewUserService() *UserService {
	return &UserService{
		users:  make(map[int]User),
		nextID: 1,
	}
}

// Create создает нового пользователя
func (s *UserService) Create(user User) (User, error) {
	if err := user.Validate(); err != nil {
		return User{}, err
	}

	// Атомарно увеличиваем ID
	id := int(atomic.AddInt64(&s.nextID, 1) - 1)
	user.ID = id

	s.mu.Lock()
	s.users[id] = user
	s.mu.Unlock()

	return user, nil
}

// Get получает пользователя по ID
func (s *UserService) Get(id int) (User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[id]
	return user, ok
}

// List возвращает всех пользователей
func (s *UserService) List() []User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users
}

// Delete удаляет пользователя
func (s *UserService) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.users[id]; !ok {
		return false
	}

	delete(s.users, id)
	return true
}

// Server представляет HTTP сервер
type Server struct {
	service *UserService
}

// NewServer создает новый сервер
func NewServer() *Server {
	return &Server{
		service: NewUserService(),
	}
}

// respondJSON отправляет JSON ответ
func respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode JSON: %v", err)
	}
}

// respondError отправляет ответ с ошибкой
func respondError(w http.ResponseWriter, statusCode int, message string) {
	respondJSON(w, statusCode, ErrorResponse{Error: message})
}

// handleUsers обрабатывает /users
func (s *Server) handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.createUser(w, r)
	case http.MethodGet:
		s.listUsers(w, r)
	default:
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// createUser создает пользователя
func (s *Server) createUser(w http.ResponseWriter, r *http.Request) {
	// Проверяем Content-Type
	contentType := r.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		respondError(w, http.StatusUnsupportedMediaType,
			"Content-Type must be application/json")
		return
	}

	// Ограничиваем размер body (1MB)
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()

	// Декодируем JSON
	var user User
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Запрещаем неизвестные поля

	if err := decoder.Decode(&user); err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case err == io.EOF:
			respondError(w, http.StatusBadRequest, "request body is empty")
		case err.Error() == "http: request body too large":
			respondError(w, http.StatusRequestEntityTooLarge,
				"request body too large (max 1MB)")
		case strings.HasPrefix(err.Error(), "json: unknown field"):
			respondError(w, http.StatusBadRequest, err.Error())
		case strings.Contains(err.Error(), "cannot unmarshal"):
			respondError(w, http.StatusBadRequest, "invalid JSON format")
		case errors.As(err, &syntaxError):
			respondError(w, http.StatusBadRequest,
				fmt.Sprintf("malformed JSON at position %d", syntaxError.Offset))
		case errors.As(err, &unmarshalTypeError):
			respondError(w, http.StatusBadRequest,
				fmt.Sprintf("invalid value for field %s", unmarshalTypeError.Field))
		default:
			respondError(w, http.StatusBadRequest, "invalid JSON")
		}
		return
	}

	// Создаем пользователя
	createdUser, err := s.service.Create(user)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Возвращаем созданного пользователя с 201
	respondJSON(w, http.StatusCreated, createdUser)
}

// listUsers возвращает список пользователей
func (s *Server) listUsers(w http.ResponseWriter, r *http.Request) {
	users := s.service.List()
	respondJSON(w, http.StatusOK, users)
}

// handleUser обрабатывает /user/{id}
func (s *Server) handleUser(w http.ResponseWriter, r *http.Request) {
	// Парсим ID из URL
	idStr := strings.TrimPrefix(r.URL.Path, "/user/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.getUser(w, r, id)
	case http.MethodDelete:
		s.deleteUser(w, r, id)
	default:
		respondError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// getUser возвращает пользователя
func (s *Server) getUser(w http.ResponseWriter, r *http.Request, id int) {
	user, ok := s.service.Get(id)
	if !ok {
		respondError(w, http.StatusNotFound, "user not found")
		return
	}

	respondJSON(w, http.StatusOK, user)
}

// deleteUser удаляет пользователя
func (s *Server) deleteUser(w http.ResponseWriter, r *http.Request, id int) {
	if !s.service.Delete(id) {
		respondError(w, http.StatusNotFound, "user not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// loggingMiddleware логирует запросы
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next(w, r)
	}
}

func main() {
	server := NewServer()

	// Регистрируем handlers с middleware
	http.HandleFunc("/users", loggingMiddleware(server.handleUsers))
	http.HandleFunc("/user/", loggingMiddleware(server.handleUser))

	fmt.Println("Server starting on :8080")
	fmt.Println("Try:")
	fmt.Println("  curl -X POST http://localhost:8080/users -H 'Content-Type: application/json' -d '{\"name\":\"John\",\"email\":\"john@example.com\",\"age\":30}'")
	fmt.Println("  curl http://localhost:8080/users")
	fmt.Println("  curl http://localhost:8080/user/1")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// ПОДРОБНОЕ ОБЪЯСНЕНИЕ ПРОБЛЕМ:
//
// ПРОБЛЕМА 1: DoS через большой JSON
//
// Атака:
//   curl -X POST http://localhost:8080/users \
//     --data-binary @10GB.json
//
// Без ограничения:
//   io.ReadAll(r.Body) → читает весь body в память
//   10GB body → 10GB в памяти
//   Server Out of Memory → crash ☠️
//
// С ограничением:
//   r.Body = http.MaxBytesReader(w, r.Body, 1<<20)  // 1MB
//   10GB body → читается только 1MB
//   decoder.Decode() → error: request body too large
//   → возвращаем 413 Request Entity Too Large ✓
//
// ПРОБЛЕМА 2: Отсутствие валидации
//
// Без валидации:
//   POST {"name":"","email":"invalid","age":-100}
//   → создается пользователь с невалидными данными
//   → в базе данных мусор
//
// С валидацией:
//   if name == "" → error: "name is required"
//   if !strings.Contains(email, "@") → error: "invalid email"
//   if age < 0 || age > 150 → error: "age must be between 0 and 150"
//   → 400 Bad Request
//
// ПРОБЛЕМА 3: Race condition в map
//
// Без синхронизации:
//   Request 1: users[1] = User{...}  ← запись
//   Request 2: user := users[1]      ← чтение
//   ↑ RACE! panic: concurrent map read and map write
//
// С RWMutex:
//   Request 1: Lock() → users[1] = ... → Unlock()
//   Request 2: RLock() → user := users[1] → RUnlock()
//   ↑ Безопасно! Мьютекс синхронизирует доступ
//
// ПРОБЛЕМА 4: Race condition в nextID
//
// Без атомарности:
//   Request 1: id := nextID; nextID++  // id = 1
//   Request 2: id := nextID; nextID++  // id = 1 (тоже!)
//   ↑ Два пользователя с одинаковым ID!
//
// С atomic:
//   Request 1: id := atomic.AddInt64(&nextID, 1)  // id = 1
//   Request 2: id := atomic.AddInt64(&nextID, 1)  // id = 2
//   ↑ Уникальные ID гарантированы
//
// ПРОБЛЕМА 5: Неправильные HTTP статусы
//
// Без правильных статусов:
//   POST /users → 200 OK (должно быть 201 Created)
//   GET /user/999 → 200 OK с пустым телом (должно быть 404)
//   Невалидный JSON → 200 OK (должно быть 400)
//   Атакующий не знает что пошло не так
//
// С правильными статусами:
//   POST /users → 201 Created
//   GET /user/999 → 404 Not Found
//   Невалидный JSON → 400 Bad Request
//   PUT /users → 405 Method Not Allowed
//   Клиент четко понимает результат
//
// ПРОБЛЕМА 6: Content-Type не проверяется
//
// Без проверки:
//   curl -X POST http://localhost:8080/users \
//     -H 'Content-Type: text/plain' \
//     -d 'not json at all'
//   → json.Unmarshal не распарсит
//   → создастся пустой User{}
//
// С проверкой:
//   if !strings.Contains(contentType, "application/json") {
//     return 415 Unsupported Media Type
//   }
//
// ПРОБЛЕМА 7: DisallowUnknownFields
//
// Без DisallowUnknownFields:
//   POST {"name":"John","hacker":"payload"}
//   → декодируется как User{Name:"John"}
//   → поле "hacker" игнорируется
//   → атакующий может не заметить что поле не работает
//
// С DisallowUnknownFields:
//   POST {"name":"John","hacker":"payload"}
//   → error: "json: unknown field \"hacker\""
//   → 400 Bad Request
//   → API строгий, помогает обнаружить ошибки
//
// ПРИМЕРЫ ЗАПРОСОВ:
//
// Создание пользователя:
//   curl -X POST http://localhost:8080/users \
//     -H 'Content-Type: application/json' \
//     -d '{"name":"John","email":"john@example.com","age":30}'
//   → 201 Created
//   → {"id":1,"name":"John","email":"john@example.com","age":30}
//
// Невалидный email:
//   curl -X POST http://localhost:8080/users \
//     -H 'Content-Type: application/json' \
//     -d '{"name":"John","email":"invalid","age":30}'
//   → 400 Bad Request
//   → {"error":"invalid email format"}
//
// Слишком большое тело:
//   curl -X POST http://localhost:8080/users \
//     --data-binary @10MB.json
//   → 413 Request Entity Too Large
//   → {"error":"request body too large (max 1MB)"}
//
// Неправильный Content-Type:
//   curl -X POST http://localhost:8080/users \
//     -H 'Content-Type: text/plain' \
//     -d 'not json'
//   → 415 Unsupported Media Type
//   → {"error":"Content-Type must be application/json"}
//
// Несуществующий пользователь:
//   curl http://localhost:8080/user/999
//   → 404 Not Found
//   → {"error":"user not found"}
//
// ЛУЧШИЕ ПРАКТИКИ:
//
// 1. Всегда ограничивайте размер request body
// 2. Валидируйте все входные данные
// 3. Проверяйте Content-Type
// 4. Используйте правильные HTTP статусы
// 5. Синхронизируйте доступ к shared state
// 6. Обрабатывайте все ошибки
// 7. DisallowUnknownFields для строгого API
// 8. Логируйте все запросы
// 9. Используйте роутеры (gorilla/mux, chi) для production
// 10. Тестируйте с -race флагом
