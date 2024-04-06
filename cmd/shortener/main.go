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
	hostIPAddr = "localhost" // адрес сервера
	hostPort   = "8080"      // порт сервера
	baseURLArg = ""
)

var myHandler MyHandler // хэндлер для обработки запросов пользлователя

func main() {
	// указываем имя флага, значение по умолчанию и описание
	hostPortArg := flag.String("a", "8080", "server port")
	baseURLArg := flag.String("b", "", "server base URL")

	// var cfg Config
	// err := env.Parse(&cfg)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	if *baseURLArg != "" {
		hostIPAddr = "/" + *baseURLArg
	}

	hostPort = *hostPortArg

	// делаем разбор командной строки
	flag.Parse()

	// создаем хэндлер
	myHandler = MyHandler{}

	// создаем репозиторий
	myHandler.rep = repository.LocalDatabase{}

	// инициализируем структура для хранения данных в репозитории
	myHandler.rep.ReceivedURL = make(map[string]string)

	// заполняем поля хэндлера
	myHandler.hostIPAddr = hostIPAddr
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
	// mux := http.NewServeMux()

	//вместо моего mux используем chi
	routerChi := chi.NewRouter()

	// назначаем обработчики для запрсов
	// mux.HandleFunc(`/`, myHandler.ServeHTTP)
	routerChi.Get(`/`+baseURLArg+`{key}`, myHandler.processGetData)
	routerChi.Post(`/`, myHandler.processPostData)

	// запускаем сервер
	// return http.ListenAndServe(hostIPAddr+`:`+hostPort, mux)
	return http.ListenAndServe(hostIPAddr+`:`+hostPort, routerChi)
}
