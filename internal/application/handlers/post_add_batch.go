package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/cfif1982/urlshtr.git/internal/domain/links"
)

type (
	PostBatchBodyRequest struct {
		CorrelationId string `json:"correlation_id,omitempty"`
		OriginalURL   string `json:"original_url,omitempty"`
	}

	PostBatchBodyResponse struct {
		CorrelationId string `json:"correlation_id,omitempty"`
		ShortURL      string `json:"short_url,omitempty"`
	}
)

// Обрабатываем запрос на добавление ссылки в БД
func (h *Handler) PostAddBatchLink(rw http.ResponseWriter, req *http.Request) {

	// переменные для работы с marshal, Unmarshal
	var postBatchBodyRequest []PostBatchBodyRequest    // массив
	postBatchBodyResponse := []PostBatchBodyResponse{} // слайс

	// после чтения тела запроса, закрываем
	defer req.Body.Close()

	// читаем тело запроса
	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.logger.Fatal(err.Error())
	}

	// анмаршалим тело в массив структур PostBatchBodyRequest
	if err = json.Unmarshal(body, &postBatchBodyRequest); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	// перебираем получившийся массив структур
	for _, v := range postBatchBodyRequest {
		// флаг создания ссылки
		bLinkCreated := false

		// создаем переменную для хранения ссылки
		var link *links.Link

		// повторяем цикл до тех пор, пока ссылка не создастся.
		// Делаю на случай существования такого ключа
		for !bLinkCreated {

			// обращаемся к domain - создаем объект ССЫЛКА
			link, err = links.CreateLink(v.OriginalURL)
			if err != nil {
				h.logger.Fatal(err.Error())
			}

			// тут мне непонятно как быть с добавлением нескольких записей в БД одним запросом
			// можно подготовить один запрос insert, но тут проблема - нужно проверять существование сгенерированного ключа в базе данных (чтобы не было дублирования)
			// тогда после генерации ключа нужно делать запрос к базе данных и проверять - существует ли запись с таким ключом.
			// и если нет, то добавить текст запроса на добавление новой записи в один общий запрос.
			// тогда какой смысл во вставлении записей одним запросом? если всё-равно на каждую новую ссылку нажуно делать запрос к БД
			// тогда просто запишу поштучно в БД

			// обращаемся к БД - сохраняем ссылку в БД
			err = h.repo.AddLink(link)

			// если err равна links.ErrKeyAlreadyExist, то нужно повторить генерацию ссылки и сохранить ее еще раз
			// во всех других случаях заканчиваем цикл(либо успешное создание ссылки, либо другая какая ошибка)
			if err != links.ErrKeyAlreadyExist {
				bLinkCreated = true
			}
		}

		if err != nil {
			h.logger.Fatal(err.Error())
		}

		// сохраняем даные для формирования ответа сервера
		postBatchBodyResponse = append(
			postBatchBodyResponse,
			PostBatchBodyResponse{
				CorrelationId: v.CorrelationId,
				ShortURL:      link.Key(),
			})
	}

	// Устанавливаем в заголовке тип передаваемых данных
	rw.Header().Set("Content-Type", "application/json")

	// устанавливаем код 201
	rw.WriteHeader(http.StatusCreated)

	// маршалим текст ответа
	answerText, err := json.Marshal(postBatchBodyResponse)

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
