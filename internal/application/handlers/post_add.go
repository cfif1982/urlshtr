package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/cfif1982/urlshtr.git/internal/domain/links"
)

type (
	PostBody struct {
		URL    string `json:"url,omitempty"`
		Result string `json:"result,omitempty"`
	}
)

// Обрабатываем запрос на добавление ссылки в БД
func (h *Handler) PostAddLink(rw http.ResponseWriter, req *http.Request) {

	var postBody PostBody

	// после чтения тела запроса, закрываем
	defer req.Body.Close()

	// читаем тело запроса
	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.logger.Fatal(err.Error())
	}

	if err = json.Unmarshal(body, &postBody); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	// обращаемся к domain - создаем объект ССЫЛКА
	link, err := links.CreateLink(postBody.URL)
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
	postBody.Result = h.baseURL + "/" + link.Key()
	postBody.URL = ""

	answerText, err := json.Marshal(postBody)

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
