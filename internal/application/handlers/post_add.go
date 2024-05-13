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

	// флаг создания ссылки
	bLinkCreated := false

	// создаем переменную для хранения ссылки
	var link *links.Link

	// повторяем цикл до тех пор, пока ссылка не создастся.
	// Делаю на случай существования такого ключа
	for !bLinkCreated {
		// обращаемся к domain - создаем объект ССЫЛКА
		link, err = links.CreateLink(postBodyRequest.URL)
		if err != nil {
			h.logger.Fatal(err.Error())
		}

		// проверяем - есть ли такой key в
		// если ключа нет, то сохраняем ссылку в БД, иначе генерируем новую ссылку
		// если при создани возникла ошибка, то ее потом обрабатываем
		if ok := h.repo.IsKeyExist(link.Key()); !ok {
			// обращаемся к БД - сохраняем ссылку в БД
			err = h.repo.AddLink(link)

			bLinkCreated = true
		}
	}

	// Устанавливаем в заголовке тип передаваемых данных
	rw.Header().Set("Content-Type", "application/json")

	// проверяем: если ошибка links.ErrURLAlreadyExist, то выводим информацию об этом в ответе сервера
	if err != nil {
		if err == links.ErrURLAlreadyExist {
			// запрос к БД - находим ссылку по ключу
			link, _ = h.repo.GetLinkByURL(link.URL())

			// устанавливаем код 409
			rw.WriteHeader(http.StatusConflict)
		} else {
			h.logger.Fatal(err.Error())
		}
	} else {
		// устанавливаем код 201
		rw.WriteHeader(http.StatusCreated)
	}

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
