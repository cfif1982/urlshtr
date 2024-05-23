package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cfif1982/urlshtr.git/pkg/logger"

	"github.com/cfif1982/urlshtr.git/internal"
	"github.com/cfif1982/urlshtr.git/internal/application/handlers"
	"github.com/cfif1982/urlshtr.git/internal/domain/links"
	linksInfra "github.com/cfif1982/urlshtr.git/internal/infrastructure/links"

	"github.com/stretchr/testify/assert"
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

	// инициализируем логгер
	logger, _ := logger.GetLogger()

	// создаем сервер
	// Его создаем для того, чтобы можно было получить доступ к его функциям, а не для его запуска
	srv := internal.NewServer("http://localhost:8080", "http://localhost", "", "", logger)

	// создаем репозиторий
	linkRepo := linksInfra.NewLocalRepository()

	// создаем хэдлер и передаем ему нужную БД
	handler := handlers.NewHandler(linkRepo, "http://localhost:8080", logger)

	// инициализируем роутер
	routerChi := srv.InitRoutes(handler)

	// перебираем параметры для тестов
	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			// создаем тестовую запись для БД
			link, err := links.NewLink(test.dataKey, test.dataValue, 0, false)
			require.NoError(t, err)

			// Добавляем в БД тестовую запись
			err = linkRepo.AddLink(link)
			require.NoError(t, err)

			// создаем запрос методом GET
			request, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/"+test.dataKey, nil)

			request.Header.Add("Accept-Encoding", "gzip")
			request.Header.Add("Content-Type", "application/json")

			// создаем рекордер для роутера
			rec := httptest.NewRecorder()

			// выполняем запрос через роутер Chi
			routerChi.ServeHTTP(rec, request)

			// проверяем код ответа
			assert.Equal(t, test.want.code, rec.Code)

			// Проверяем заголовок ответа
			assert.Equal(t, test.want.headerValue, rec.Header().Get(test.want.headerType))

		})
	}
}
