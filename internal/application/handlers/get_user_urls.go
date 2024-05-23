package handlers

import (
	"encoding/json"
	"net/http"
)

type URLForResponse struct {
	ShortURL    string `json:"short_url,omitempty"`
	OriginalURL string `json:"original_url,omitempty"`
}

// Обрабатываем запрос на получение ссылок из БД по id пользлователя
func (h *Handler) GetUserURLS(rw http.ResponseWriter, req *http.Request) {

	arrURLForResponse := []URLForResponse{} // слайс

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

	// запрос к БД - находим ссылку по ключу
	urls, err := h.repo.GetLinksByUserID(userID)

	if err != nil {
		h.logger.Info(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// сохраняем полученные данные в нужном для вывода формате
	for _, v := range *urls {
		arrURLForResponse = append(
			arrURLForResponse,
			URLForResponse{
				ShortURL:    h.baseURL + "/" + v.Key(),
				OriginalURL: v.URL(),
			})

	}

	// Устанавливаем в заголовке тип передаваемых данных
	rw.Header().Set("Content-Type", "application/json")

	// устанавливаем код 200
	rw.WriteHeader(http.StatusOK)

	// маршалим текст ответа
	answerText, err := json.Marshal(arrURLForResponse)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// выводим ответ сервера
	_, err = rw.Write([]byte(answerText))
	if err != nil {
		h.logger.Fatal(err.Error())
	}
}
