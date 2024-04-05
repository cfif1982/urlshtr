package main

import (
	"net/http"

	"github.com/cfif1982/urlshtr.git/cmd/shortener/repository"
)

var hostIpAddr = "localhost" // адрес сервера
var hostPort = "8080"        // порт сервера

var myHandler MyHandler // хэндлер для обработки запросов пользлователя

func main() {
	// создаем хэндлер
	myHandler = MyHandler{}

	// создаем репозиторий
	myHandler.rep = repository.LocalDatabase{}

	// инициализируем структура для хранения данных в репозитории
	myHandler.rep.ReceivedURL = make(map[string]string)

	// заполняем поля хэндлера
	myHandler.hostIPAddr = hostIpAddr
	myHandler.hostPort = hostPort

	// запускем сервер
	_ = run()
	if err := run(); err != nil {
		panic(err)
	}
}

// Запуск сервера
func run() error {
	// создаем свой маршрутизатор запросов от пользователя
	mux := http.NewServeMux()

	// назначаем обработчики для запрсов
	mux.HandleFunc(`/`, myHandler.ServeHTTP)

	// запускаем сервер
	return http.ListenAndServe(hostIpAddr+`:`+hostPort, mux)
}
