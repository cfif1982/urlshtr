package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/cfif1982/urlshtr.git/internal/domain/links"
)

type (
	PostBatchBodyRequest struct {
		CorrelationID string `json:"correlation_id,omitempty"`
		OriginalURL   string `json:"original_url,omitempty"`
	}

	PostBatchBodyResponse struct {
		CorrelationID string `json:"correlation_id,omitempty"`
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

	// создаем map для хранения ссылок
	mapLinks := []*links.Link{}

	// перебираем получившийся массив структур
	for _, v := range postBatchBodyRequest {

		// флаг создания ссылки
		bLinkCreated := false

		// создаем переменную для хранения ссылки
		var link *links.Link

		// повторяем цикл до тех пор, пока ссылка не создастся.
		// Делаю на случай существования такого ключа
		for !bLinkCreated {
			// узнаем id пользователя из контекста запроса
			userID := 0
			if req.Context().Value("user_id_key") != nil {
				userID = req.Context().Value("user_id_key").(int)
			}

			// обращаемся к domain - создаем объект ССЫЛКА
			link, err = links.CreateLink(v.OriginalURL, userID)
			if err != nil {
				h.logger.Fatal(err.Error())
			}

			if err != nil {
				h.logger.Fatal(err.Error())
			}

			// проверяем - есть ли такой key в БД
			// если ключа нет, то добавляем ссылку в map, иначе генерируем новую ссылку
			if ok := h.repo.IsKeyExist(link.Key()); !ok {

				// сохраняем ссылку в map
				mapLinks = append(mapLinks, link)

				// сохраняем даные для формирования ответа сервера
				postBatchBodyResponse = append(
					postBatchBodyResponse,
					PostBatchBodyResponse{
						CorrelationID: v.CorrelationID,
						ShortURL:      h.baseURL + "/" + link.Key(),
					})

				// устанавливаем флаг создания ссылки для прекращения цикла
				bLinkCreated = true
			}
		}
	}

	// После того как сформировали массив ссылок для добавлеия в БД, добавляем всё дним запросом
	err = h.repo.AddLinkBatch(mapLinks)
	if err != nil {
		h.logger.Fatal(err.Error())
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
