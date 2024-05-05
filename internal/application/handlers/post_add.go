package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/cfif1982/urlshtr.git/internal/domain/links"
)

type (
	PostBodyRequest struct {
		URL string `json:"url,omitempty"`
	}

	PostBodyResponse struct {
		Result string `json:"result,omitempty"`
	}
)

// Обрабатываем запрос на добавление ссылки в БД
func (h *Handler) PostAddLink(rw http.ResponseWriter, req *http.Request) {

	var postBodyRequest PostBodyRequest

	// после чтения тела запроса, закрываем
	defer req.Body.Close()

	// читаем тело запроса
	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.logger.Fatal(err.Error())
	}

	if err = json.Unmarshal(body, &postBodyRequest); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	// обращаемся к domain - создаем объект ССЫЛКА
	link, err := links.CreateLink(postBodyRequest.URL)
	if err != nil {
		h.logger.Fatal(err.Error())
	}

	// обращаемся к БД - сохраняем ссылку в БД
	err = h.repo.AddLink(link)

	if err != nil {
		h.logger.Fatal(err.Error())
	}

	// Устанавливаем в заголовке тип передаваемых данных
	rw.Header().Set("Content-Type", "application/json")

	// устанавливаем код 201
	rw.WriteHeader(http.StatusCreated)

	// формируем текст ответа сервера
	postBodyResponse := PostBodyResponse{
		Result: h.baseURL + "/" + link.Key(),
	}

	answerText, err := json.Marshal(postBodyResponse)

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
