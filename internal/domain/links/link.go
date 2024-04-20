package links

import (
	"errors"

	"github.com/google/uuid"
)

// список возможных шибок
var (
	ErrLinkNotFound    = errors.New("link not found")
	ErrKeyAlreadyExist = errors.New("key already exist")
)

// структура для хранения объекта ССЫЛКА
type Link struct {
	key string
	url string
}

// создаем новый объект ССЫЛКА
// нужна для использвания в других пакетах
func NewLink(key string, url string) (*Link, error) {
	return &Link{
		key: key,
		url: url,
	}, nil
}

// Создаем новую ССЫЛКУ
func CreateLink(url string) (*Link, error) {

	key := generateKey()

	return NewLink(key, url)
}

// генерируем key
func generateKey() string {

	// генерируем случайный код типа string
	uuid := uuid.NewString()[:8]

	return uuid
}

// возвращщаем поле key
func (l *Link) Key() string {
	return l.key
}

// возвращщаем поле URL
func (l *Link) URL() string {
	return l.url
}
