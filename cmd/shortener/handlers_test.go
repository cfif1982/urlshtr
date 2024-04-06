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
				response:    serverBaseURL + "/",
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
		myHandler.serverAddress = serverAddress
		myHandler.serverBaseURL = serverBaseURL

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

// func testRequest(t *testing.T, ts *httptest.Server, key string) *http.Response {
// 	req, err := http.NewRequest(http.MethodGet, ts.URL+"/"+key, nil)
// 	require.NoError(t, err)

// 	resp, err := ts.Client().Do(req)
// 	require.NoError(t, err)
// 	defer resp.Body.Close()

// 	return resp
// }
/*
func TestProcessGetData(t *testing.T) {
	myHandler := MyHandler{}

	// создаем репозиторий
	myHandler.rep = repository.LocalDatabase{}

	// инициализируем структура для хранения данных в репозитории
	myHandler.rep.ReceivedURL = make(map[string]string)

	// заполняем поля хэндлера
	myHandler.hostIPAddr = hostIPAddr
	myHandler.hostPort = hostPort

	routerChi := chi.NewRouter()
	routerChi.Get(`/{key}`, myHandler.processGetData)

	// testBody := `testBody: /` + baseURLArg + `{key}`
	// fmt.Printf("%v", testBody)

	ts := httptest.NewServer(routerChi)
	defer ts.Close()

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

		// Добавляем в БД тестовую запись
		myHandler.rep.ReceivedURL[test.dataKey] = test.dataValue

		// testBody := `testBody: /` + baseURLArg + `{key}`
		fmt.Printf("%v", myHandler.rep.ReceivedURL[test.dataKey])

		// создаем запрос методом GET
		response := testRequest(t, ts, test.dataKey)

		// // создаём новый Recorder
		// w := httptest.NewRecorder()
		// myHandler.processGetData(w, request)

		// res := w.Result()

		// body := "Header: \r\n"
		// for k, v := range res.Header {
		// 	body += fmt.Sprintf("%s: %v\r\n", k, v)
		// }

		// fmt.Printf("%v", body)

		// // получаем и проверяем тело запроса
		// defer res.Body.Close()

		// проверяем код ответа
		assert.Equal(t, test.want.code, response.StatusCode)

		// Проверяем заголовок ответа
		assert.Equal(t, test.want.headerValue, response.Header.Get(test.want.headerType))
	}
}
*/
