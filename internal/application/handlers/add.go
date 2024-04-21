package handlers

import (
	"io"

	"net/http"

	"github.com/cfif1982/urlshtr.git/internal/domain/links"
)

// Обрабатываем запрос на добавление ссылки в БД
func (h *Handler) AddLink(rw http.ResponseWriter, req *http.Request) {

	// SugarLogger.Infow(
	// 	"POST request received",
	// )

	// после чтения тела запроса, закрываем
	defer req.Body.Close()

	// читаем тело запроса
	// body, err := io.ReadAll(req.Body)
	body, _ := io.ReadAll(req.Body)
	// if err != nil {
	// 	SugarLogger.Fatal(err)
	// }

	// обращаемся к domain - создаем объект ССЫЛКА
	// link, err := links.CreateLink(string(body))
	link, _ := links.CreateLink(string(body))
	// if err != nil {
	// 	SugarLogger.Fatal(err)
	// }

	// обращаемся к БД - сохраняем ссылку в БД
	// err = h.repo.AddLink(link)
	_ = h.repo.AddLink(link)

	// if err != nil {
	// 	SugarLogger.Fatal(err)
	// }

	// Устанавливаем в заголовке тип передаваемых данных
	rw.Header().Set("Content-Type", "text/plain")

	// устанавливаем код 201
	rw.WriteHeader(http.StatusCreated)

	// формируем текст ответа сервера
	answerText := h.baseURL + "/" + link.Key()

	// выводим ответ сервера
	// _, err = rw.Write([]byte(answerText))
	_, _ = rw.Write([]byte(answerText))
	// if err != nil {
	// 	SugarLogger.Fatal(err)
	// }

}
