package main

import (
	"flag"
	"net/http"

	"github.com/cfif1982/urlshtr.git/cmd/shortener/repository"
	"github.com/go-chi/chi/v5"
)

// type Config struct {
// 	serverAddress string `env:"SERVER_ADDRESS"`
// 	baseURL       string `env:"BASE_URL"`
// }

var (
	serverAddress string // адрес сервера
	serverBaseURL string // порт сервера
)

var myHandler MyHandler // хэндлер для обработки запросов пользлователя

func main() {
	// указываем имя флага, значение по умолчанию и описание
	serverAddressArg := flag.String("a", "localhost:8080", "server address ")
	serverBaseURLArg := flag.String("b", "http://localhost:8080", "server base URL")

	// делаем разбор командной строки
	flag.Parse()

	// var cfg Config
	// err := env.Parse(&cfg)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	serverAddress = *serverAddressArg
	serverBaseURL = *serverBaseURLArg

	// создаем хэндлер
	myHandler = MyHandler{}

	// создаем репозиторий
	myHandler.rep = repository.LocalDatabase{}

	// инициализируем структура для хранения данных в репозитории
	myHandler.rep.ReceivedURL = make(map[string]string)

	// заполняем поля хэндлера
	myHandler.serverAddress = serverAddress
	myHandler.serverBaseURL = serverBaseURL

	// запускем сервер
	_ = run()
	if err := run(); err != nil {
		panic(err)
	}
}

// Запуск сервера
func run() error {
	// создаем свой маршрутизатор запросов от пользователя
	// mux := http.NewServeMux()

	//вместо моего mux используем chi
	routerChi := chi.NewRouter()

	// назначаем обработчики для запрсов
	// mux.HandleFunc(`/`, myHandler.ServeHTTP)
	routerChi.Get(`/{key}`, myHandler.processGetData)
	routerChi.Post(`/`, myHandler.processPostData)

	// запускаем сервер
	// return http.ListenAndServe(hostIPAddr+`:`+hostPort, mux)
	return http.ListenAndServe(serverAddress, routerChi)
}
