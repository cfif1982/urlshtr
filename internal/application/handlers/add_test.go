package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cfif1982/urlshtr.git/internal"
	"github.com/cfif1982/urlshtr.git/internal/application/handlers"
	linksInfra "github.com/cfif1982/urlshtr.git/internal/infrastructure/links"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddLink(t *testing.T) {
	type want struct {
		code        int
		response    string
		headerType  string
		headerValue string
	}
	tests := []struct {
		name        string
		requestBody string
		want        want
	}{
		{
			name:        "add link test #1",
			requestBody: "https://practicum.yandex.ru/",
			want: want{
				code:        http.StatusCreated,
				response:    "http://localhost:8080/",
				headerType:  "Content-Type",
				headerValue: "text/plain",
			},
		},
	}

	// создаем сервер
	// Его создаем для того, чтобы можно было получить доступ к его функциям, а не для его запуска
	// srv := new(internal.Server)
	srv := internal.NewServer("http://localhost:8080", "http://localhost")

	// устанавливаем данные из флагов и переменных среды
	// srv.SetServerAddress("http://localhost:8080")
	// srv.SetServerBaseURL("http://localhost")

	// создаем репозиторий
	linkRepo := linksInfra.NewLocalRepository()

	// создаем хэдлер и передаем ему нужную БД
	handler := handlers.NewHandler(linkRepo, srv.GetServerAddress())
	//********************************************************

	// инициализируем роутер
	routerChi := srv.InitRoutes(handler)

	// создаем тестовый сервер
	ts := httptest.NewServer(routerChi)

	// перебираем параметры для тестов
	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			// готовим текст для передачи  в тело запроса
			body := strings.NewReader(test.requestBody)

			// создаем запрос методом POST
			request, _ := http.NewRequest(http.MethodPost, ts.URL+"/", body)

			// выполняем запрос
			resp, err := ts.Client().Do(request)
			require.NoError(t, err)

			// проверяем код ответа
			assert.Equal(t, test.want.code, resp.StatusCode)

			// получаем тело запроса
			defer resp.Body.Close()
			resBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			// в теле ответа должна появиться ссылка - находим в ней key
			testedKey := string(resBody)[len(test.want.response):]

			// проверяем тело запроса
			assert.Equal(t, test.want.response+testedKey, string(resBody))

			// Проверяем заголовок ответа
			assert.Equal(t, test.want.headerValue, resp.Header.Get(test.want.headerType))
		})
	}
}
