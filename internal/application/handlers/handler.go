package handlers

import (
	"github.com/cfif1982/urlshtr.git/pkg/log"

	"github.com/cfif1982/urlshtr.git/internal/domain/links"
)

// Интерфейс репозитория
type RepositoryInterface interface {

	// Добавить ссылку в БД
	AddLink(link *links.Link) error

	// Найти ссылку в БД по key
	GetLinkByKey(key string) (*links.Link, error)

	// узнаем - доступна ли БД
	Ping() error
}

// структура хэндлера
type Handler struct {
	repo    RepositoryInterface
	baseURL string
	logger  *log.Logger
}

// создаем новый хэндлер
func NewHandler(repo RepositoryInterface, base string, logger *log.Logger) *Handler {
	return &Handler{
		repo:    repo,
		baseURL: base,
		logger:  logger,
	}
}
