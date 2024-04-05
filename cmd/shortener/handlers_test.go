package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cfif1982/urlshtr.git/cmd/shortener/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcessPostData(t *testing.T) {
	type want struct {
		code        int
		response    string
		headerType  string
		headerValue string
	}
	tests := []struct {
		name        string
		requestData string
		want        want
	}{
		{
			name:        "post request test #1",
			requestData: "https://practicum.yandex.ru/",
			want: want{
				code:        http.StatusCreated,
				response:    "http://" + hostIPAddr + ":" + hostPort + "/",
				headerType:  "Content-Type",
				headerValue: "text/plain",
			},
		},
	}
	for _, test := range tests {
		myHandler := MyHandler{}

		// создаем репозиторий
		myHandler.rep = repository.LocalDatabase{}

		// инициализируем структура для хранения данных в репозитории
		myHandler.rep.ReceivedURL = make(map[string]string)

		// заполняем поля хэндлера
		myHandler.hostIPAddr = hostIPAddr
		myHandler.hostPort = hostPort

		t.Run(test.name, func(t *testing.T) {
			// создаем запрос методом POST
			body := strings.NewReader(test.requestData)
			request := httptest.NewRequest(http.MethodPost, "/", body)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			myHandler.processPostData(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)

			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody)[:len(test.want.response)])

			// Проверяем заголовок ответа
			assert.Equal(t, test.want.headerValue, res.Header.Get(test.want.headerType))
		})
	}
}

func TestProcessGetData(t *testing.T) {
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
	for _, test := range tests {
		myHandler := MyHandler{}

		// создаем репозиторий
		myHandler.rep = repository.LocalDatabase{}

		// инициализируем структура для хранения данных в репозитории
		myHandler.rep.ReceivedURL = make(map[string]string)

		// заполняем поля хэндлера
		myHandler.hostIPAddr = hostIPAddr
		myHandler.hostPort = hostPort

		t.Run(test.name, func(t *testing.T) {
			// Добавляем в БД тестовую запись
			myHandler.rep.ReceivedURL[test.dataKey] = test.dataValue

			// создаем запрос методом GET
			request := httptest.NewRequest(http.MethodGet, "/"+test.dataKey, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			myHandler.processGetData(w, request)

			res := w.Result()

			// получаем и проверяем тело запроса
			defer res.Body.Close()

			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)

			// Проверяем заголовок ответа
			assert.Equal(t, test.want.headerValue, res.Header.Get(test.want.headerType))
		})
	}
}
