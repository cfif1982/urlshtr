package handlers

import (
	"github.com/cfif1982/urlshtr.git/pkg/log"

	"github.com/cfif1982/urlshtr.git/internal/domain/links"
)

// Интерфейс репозитория
type RepositoryInterface interface {

	// узнаем - есть ли уже запись с данным ключом
	IsKeyExist(key string) bool

	// Добавить ссылку в БД
	AddLink(link *links.Link) error

	// Добавить массив ссылок в БД
	AddLinkBatch(links []*links.Link) error

	// Найти ссылку в БД по key
	GetLinkByKey(key string) (*links.Link, error)

	// Найти ссылку в БД по url
	GetLinkByURL(key string) (*links.Link, error)

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
