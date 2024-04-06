package main

import (
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/cfif1982/urlshtr.git/cmd/shortener/repository"
	"github.com/google/uuid"
)

// Структура для хранения хэндлера
type MyHandler struct {
	rep           repository.LocalDatabase
	serverAddress string
	serverBaseURL string
}

// Обработчик запросов от польлзователя
func (h MyHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	// Если данные переданы методом POST
	if req.Method == http.MethodPost {
		h.processPostData(res, req)
	} else
	// Если данные переданы методом GET
	if req.Method == http.MethodGet {
		h.processGetData(res, req)
	} else
	// В других случаях ошибка
	{
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

// Обрабатываем данные полученные методом POST
func (h MyHandler) processPostData(res http.ResponseWriter, req *http.Request) {

	// генерируем случайный код типа string
	uuid := uuid.NewString()[:8]

	// после чтения тела запроса, закрываем
	defer req.Body.Close()

	// читаем тело запроса
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// сохраняем полученные данные
	h.rep.SaveURL(uuid, string(body))

	// формируем текст ответа сервера
	answerText := h.serverBaseURL + "/" + uuid

	// Устанавливаем в заголовке тип передаваемых данных
	res.Header().Set("Content-Type", "text/plain")

	// устанавливаем код 201
	res.WriteHeader(http.StatusCreated)

	// выводим ответ сервера
	res.Write([]byte(answerText))
}

// Обрабатываем данные полученные методом GET
func (h MyHandler) processGetData(res http.ResponseWriter, req *http.Request) {

	// узнаем данные из полученной адресной строки
	// key := req.URL.Path[1:]

	// Узнаем key с помощью chi
	key := chi.URLParam(req, "key")

	// по ключу находим данные в БД
	value := h.rep.GetURL(key)

	// Устанавливаем заголовок ответа
	res.Header().Set("Location", value)

	// устанавливаем код ответа 307
	res.WriteHeader(http.StatusTemporaryRedirect)
}
