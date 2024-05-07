package handlers

import (
	"net/http"
)

// Обрабатываем запрос на получение ссылки из БД по ключу
func (h *Handler) Ping(rw http.ResponseWriter, req *http.Request) {

	err := h.repo.Ping()

	if err == nil {
		// устанавливаем код ответа 200
		rw.WriteHeader(http.StatusOK)
	} else {

		// устанавливаем код ответа 500
		rw.WriteHeader(http.StatusInternalServerError)

		h.logger.Info(err.Error())
	}
}
