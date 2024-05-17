package handlers

import (
	"io"

	"net/http"

	"github.com/cfif1982/urlshtr.git/internal/domain/links"
)

// Обрабатываем запрос на добавление ссылки в БД
func (h *Handler) AddLink(rw http.ResponseWriter, req *http.Request) {

	// после чтения тела запроса, закрываем
	defer req.Body.Close()

	// читаем тело запроса
	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.logger.Fatal(err.Error())
	}

	// флаг создания ссылки
	bLinkCreated := false

	// создаем переменную для хранения ссылки
	var link *links.Link

	// повторяем цикл до тех пор, пока ссылка не создастся.
	// Делаю на случай существования такого ключа
	for !bLinkCreated {
		// узнаем id пользователя из контекста запроса
		userID := 0
		if req.Context().Value("userID") != nil {
			userID = req.Context().Value("userID").(int)
		}

		// обращаемся к domain - создаем объект ССЫЛКА
		link, err = links.CreateLink(string(body), userID)
		if err != nil {
			h.logger.Fatal(err.Error())
		}

		// проверяем - есть ли такой key в
		// если ключа нет, то сохраняем ссылку в БД, иначе генерируем новую ссылку
		// если при создани возникла ошибка, то ее потом обрабатываем
		ok := h.repo.IsKeyExist(link.Key())

		if !ok {
			// обращаемся к БД - сохраняем ссылку в БД
			err = h.repo.AddLink(link)

			bLinkCreated = true
		}
	}

	// Устанавливаем в заголовке тип передаваемых данных
	rw.Header().Set("Content-Type", "text/plain")

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
	answerText := h.baseURL + "/" + link.Key()

	// выводим ответ сервера
	_, err = rw.Write([]byte(answerText))
	if err != nil {
		h.logger.Fatal(err.Error())
	}
}
