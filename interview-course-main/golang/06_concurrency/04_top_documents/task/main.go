// Задача: Top K документов для пользователя
//
// Дана функция Scorer, которая оценивает релевантность документа для пользователя.
// Чем выше скор, тем релевантнее документ.
//
// Требования:
// Реализовать сервис с двумя методами:
//
// 1. AddDocument - добавить/обновить документ в систему
// 2. GetTopDocuments - получить топ K документов для пользователя
//    отсортированных по убыванию скора
//
// Особенности:
// - Документ может обновляться (тот же DocID, но новый контент)
// - GetTopDocuments должен работать быстро (используется в realtime)
// - Скор документа может меняться при его обновлении
// - Сервис работает в многопоточной среде (конкурентный доступ)
// - Может быть много документов (миллионы) и пользователей (тысячи)
//
// Пример:
// scorer.GetScore(doc1, user1) = 95
// scorer.GetScore(doc2, user1) = 80
// scorer.GetScore(doc3, user1) = 70
//
// GetTopDocuments(user1, limit=2) → [doc1, doc2]

package main

type Document struct {
	DocID   string
	Content string
	// Дополнительные поля по необходимости
}

type User struct {
	UserID string
	// Дополнительные поля по необходимости
}

// Scorer интерфейс для оценки релевантности документа
type Scorer interface {
	GetScore(doc Document, user User) int
}

// DocumentService интерфейс сервиса
type DocumentService interface {
	// AddDocument добавляет или обновляет документ
	AddDocument(doc Document) error

	// GetTopDocuments возвращает топ K документов для пользователя
	// Документы отсортированы по убыванию скора
	GetTopDocuments(user User, limit int) ([]Document, error)
}

// TODO: Реализуйте структуру и методы интерфейса DocumentService
