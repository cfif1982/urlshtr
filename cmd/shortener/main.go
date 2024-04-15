package main

import (
	"flag"
	"log"

	"github.com/caarlos0/env"
	"github.com/cfif1982/urlshtr.git/internal"
)

// храним значения переменных среды
type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}

// глобальные переменные для настройки
// var (
// 	serverAddress string // адрес сервера
// 	serverBaseURL string // порт сервера
// )

func main() {

	// указываем имя флага, значение по умолчанию и описание
	serverAddressArg := flag.String("a", "localhost:8080", "server address ")
	serverBaseURLArg := flag.String("b", "http://localhost:8080", "server base URL")

	// делаем разбор командной строки
	flag.Parse()

	var cfg Config

	// парсим переменные среды в структуру
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	// адрес сервера из флага
	serverAddress := *serverAddressArg

	// базовый URL из флага
	serverBaseURL := *serverBaseURLArg

	if cfg.ServerAddress != "" {
		serverAddress = cfg.ServerAddress
	}

	if cfg.BaseURL != "" {
		serverBaseURL = cfg.BaseURL
	}

	// создаем сервер
	srv := new(internal.Server)

	// устанавливаем данные из флагов и переменных среды
	srv.SetServerAddress(serverAddress)
	srv.SetServerBaseURL(serverBaseURL)

	// запускаем сервер
	if err := srv.Run(srv.GetServerAddress()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}

}
