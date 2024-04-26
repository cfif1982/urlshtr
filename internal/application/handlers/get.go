package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Обрабатываем запрос на получение ссылки из БД по ключу
func (h *Handler) GetLinkByKey(rw http.ResponseWriter, req *http.Request) {
	// разбираем url запроса и ищем поле key
	key := chi.URLParam(req, "key")

	// запрос к БД - находим ссылку по ключу
	url, err := h.repo.GetLinkByKey(key)

	// Если запись не найдена в БД
	if err != nil {
		h.logger.Info("link not found")
		rw.WriteHeader(http.StatusInternalServerError)
	} else {

		// Устанавливаем заголовок ответа
		rw.Header().Set("Location", url.URL())

		// устанавливаем код ответа 307
		rw.WriteHeader(http.StatusTemporaryRedirect)
	}
}
