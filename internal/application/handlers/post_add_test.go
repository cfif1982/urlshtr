package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cfif1982/urlshtr.git/pkg/log"

	"github.com/cfif1982/urlshtr.git/internal"
	"github.com/cfif1982/urlshtr.git/internal/application/handlers"
	linksInfra "github.com/cfif1982/urlshtr.git/internal/infrastructure/links"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostAddLink(t *testing.T) {
	type want struct {
		code          int
		responseStart string
		responseEnd   string
		headerType    string
		headerValue   string
	}
	tests := []struct {
		name        string
		requestBody string
		want        want
	}{
		{
			name:        "add link test #1",
			requestBody: `{"url":"https://practicum.yandex.ru"}`,
			want: want{
				code:          http.StatusCreated,
				responseStart: `{"result":"http://localhost:8080/`,
				responseEnd:   `"}`,
				headerType:    "Content-Type",
				headerValue:   "application/json",
			},
		},
	}

	// инициализируем логгер
	logger, _ := log.GetLogger()

	// создаем сервер
	// Его создаем для того, чтобы можно было получить доступ к его функциям, а не для его запуска

	srv := internal.NewServer("http://localhost:8080", "http://localhost", "", "", logger)


	// создаем репозиторий
	linkRepo := linksInfra.NewLocalRepository()

	// создаем хэдлер и передаем ему нужную БД
	handler := handlers.NewHandler(linkRepo, srv.GetServerAddress(), logger)
	//********************************************************

	// инициализируем роутер
	routerChi := srv.InitRoutes(handler)

	// перебираем параметры для тестов
	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			// готовим текст для передачи  в тело запроса
			body := strings.NewReader(test.requestBody)

			// создаем запрос методом POST
			request, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/api/shorten", body)

			// создаем рекордер для роутера
			rec := httptest.NewRecorder()

			// выполняем запрос через роутер Chi
			routerChi.ServeHTTP(rec, request)

			// проверяем код ответа
			assert.Equal(t, test.want.code, rec.Code)

			// получаем тело запроса
			defer request.Body.Close()
			resBody, err := io.ReadAll(rec.Body)
			require.NoError(t, err)

			// в теле ответа должна появиться ссылка - находим в ней key
			testedKey := string(resBody)[len(test.want.responseStart):]

			// проверяем тело запроса
			assert.Equal(t, test.want.responseStart+testedKey, string(resBody))

			// Проверяем заголовок ответа
			assert.Equal(t, test.want.headerValue, rec.Header().Get(test.want.headerType))
		})
	}
}
