package handlers

import (
	"encoding/json"
	"io"
	"net/http"
)

// Обрабатываем запрос на получение ссылок из БД по id пользлователя
func (h *Handler) DeleteUserURLS(rw http.ResponseWriter, req *http.Request) {

	// узнаем id пользователя из контекста запроса
	userID := 0
	if req.Context().Value(KeyUserID) != nil {
		userID = req.Context().Value(KeyUserID).(int)
	}

	// Если пользователь не авторизован, то выдаем собщение об этом
	if userID == 0 {
		http.Error(rw, "", http.StatusUnauthorized)
		return
	}

	// после чтения тела запроса, закрываем
	defer req.Body.Close()

	// читаем тело запроса
	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.logger.Fatal(err.Error())
	}

	var keyStings []string

	// анмаршалим тело в массив строк
	if err = json.Unmarshal(body, &keyStings); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	// запрос к БД - находим ссылку по ключу
	err = h.repo.ChangeDeletedFlagByUserID(userID, keyStings)

	if err != nil {
		h.logger.Info(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Устанавливаем в заголовке тип передаваемых данных
	rw.Header().Set("Content-Type", "application/json")

	// устанавливаем код 200
	rw.WriteHeader(http.StatusAccepted)
}
