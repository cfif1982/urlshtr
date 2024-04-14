package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/caarlos0/env"
	"github.com/cfif1982/urlshtr.git/internal"
)

// храним значения переменных среды
type Config struct {
	serverAddress string `env:"SERVER_ADDRESS"`
	baseURL       string `env:"BASE_URL"`
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

	// // если флаг не передавали, то смотрим переменную окруения
	// if !isFlagPassed("a") {
	// 	if cfg.serverAddress != "" {
	// 		serverAddress = cfg.serverAddress
	// 	}
	// }

	// // если флаг не передавали, то смотрим переменную окруения
	// if !isFlagPassed("b") {
	// 	if cfg.baseURL != "" {
	// 		serverBaseURL = cfg.baseURL
	// 	}
	// }

	if cfg.serverAddress != "" {
		serverAddress = cfg.serverAddress
	}

	if cfg.baseURL != "" {
		serverBaseURL = cfg.baseURL
	}

	// создаем сервер
	srv := new(internal.Server)

	// устанавливаем данные из флагов и переменных среды
	srv.SetServerAddress(serverAddress)
	srv.SetServerBaseURL(serverBaseURL)

	fmt.Print("serverAddress: ")
	fmt.Println(serverAddress)

	// запускаем сервер
	if err := srv.Run(srv.ServerAddress()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}

}

// проверяем передавали флаг или нет
func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
