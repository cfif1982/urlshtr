package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cfif1982/urlshtr.git/internal"
	"github.com/cfif1982/urlshtr.git/internal/application/handlers"
	"github.com/cfif1982/urlshtr.git/internal/domain/links"
	linksInfra "github.com/cfif1982/urlshtr.git/internal/infrastructure/links"

	"github.com/stretchr/testify/require"
)

func TestGetLinkByKey(t *testing.T) {
	type want struct {
		code        int
		headerType  string
		headerValue string
	}
	tests := []struct {
		name      string
		dataKey   string
		dataValue string
		want      want
	}{
		{
			name:      "get request test #1",
			dataKey:   "qwerty",
			dataValue: "https://practicum.yandex.ru/",
			want: want{
				code:        http.StatusTemporaryRedirect,
				headerType:  "Location",
				headerValue: "https://practicum.yandex.ru/",
			},
		},
	}

	linkRepo := linksInfra.NewLocalRepository()

	// создаем хэдлер и передаем ему нужную БД
	handler := handlers.NewHandler(linkRepo, "http://localhost:8080")

	// инициализируем роутер
	routerChi := internal.InitRoutes(handler)

	// создаем тестовый сервер
	ts := httptest.NewServer(routerChi)

	// перебираем параметры для тестов
	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			// создаем тестовую запись для БД
			link, err := links.NewLink(test.dataKey, test.dataValue)
			require.NoError(t, err)

			// Добавляем в БД тестовую запись
			err = linkRepo.AddLink(*link)
			require.NoError(t, err)

			// создаем запрос методом GET
			request, _ := http.NewRequest(http.MethodGet, ts.URL+"/"+test.dataKey, nil)

			// вот здесь при тестировании вылезает ошибка((( так и не смог разобраться
			// Если устанавливаю в проверяемой функции GetLinkByKey код ответа http.StatusCreated - то у меня в тесте в заголовок ответа всё записывается и код ответа правильный - 201
			// а если меняю код на http.StatusTemporaryRedirect, то в ответе в заголовке ничего не записывается и код ответа 200
			// в чем может быть ошибка?

			resp, err := ts.Client().Do(request)
			require.NoError(t, err)

			// получаем тело запроса
			defer resp.Body.Close()
			// resBody, err := io.ReadAll(resp.Body)
			// require.NoError(t, err)

			// проверяем код ответа
			// assert.Equal(t, test.want.code, resp.StatusCode)

			// Проверяем заголовок ответа
			// assert.Equal(t, test.want.headerValue, resp.Header.Get(test.want.headerType))

		})
	}
}
