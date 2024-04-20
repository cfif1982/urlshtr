package handlers

import "github.com/cfif1982/urlshtr.git/internal/domain/links"

// Интерфейс репозитория
type RepositoryInterface interface {

	// Добавить ссылку в БД
	AddLink(link *links.Link) error

	// Найти ссылку в БД по key
	GetLinkByKey(key string) (*links.Link, error)
}

// структура хэндлера
type Handler struct {
	repo    RepositoryInterface
	baseURL string
}

// создаем новый хэндлер
func NewHandler(repo RepositoryInterface, base string) *Handler {
	return &Handler{
		repo:    repo,
		baseURL: base,
	}
}
