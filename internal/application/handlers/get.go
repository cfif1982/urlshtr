package handlers

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Обрабатываем запрос на получение ссылки из БД по ключу
func (h *Handler) GetLinkByKey(rw http.ResponseWriter, req *http.Request) {

	// разбираем url запроса и ищем поле key
	key := chi.URLParam(req, "key")

	// запрос к БД - находим ссылку по ключу
	url, err := h.repo.GetLinkByKey(key)

	if err != nil {
		log.Fatalln(err)
	}

	// Устанавливаем заголовок ответа
	rw.Header().Set("Location", url.URL())

	// вот здесь при тестировании вылезает ошибка((( так и не смог разобраться
	// Если устанавливаю код http.StatusCreated - то у меня в тесте в заголовок ответа всё записывается и код ответа правильный - 201
	// а если меняю код на http.StatusTemporaryRedirect, то в ответе в заголовке ничего не записывается и код ответа 200
	// в чем может быть ошибка?
	// устанавливаем код ответа 307
	rw.WriteHeader(http.StatusTemporaryRedirect)
	// rw.WriteHeader(http.StatusCreated)

}
