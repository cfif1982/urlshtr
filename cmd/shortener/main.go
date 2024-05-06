package main

import (
	"flag"

	"github.com/caarlos0/env"
	"github.com/cfif1982/urlshtr.git/internal"
	"github.com/cfif1982/urlshtr.git/pkg/log"
)

// храним значения переменных среды
type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`

	DatabaseDSN     string `env:"DATABASE_DSN"`


}

func main() {

	// инициализируем логгер
	logger, err := log.GetLogger()

	if err != nil {
		logger.Fatal("cannot initialize zap")
	}

	// выводим сообщенеи об успешной инициализации логгера
	logger.Info("logger zap initialization: SUCCESS")

	// указываем имя флага, значение по умолчанию и описание
	serverAddressArg := flag.String("a", "localhost:8080", "server address ")
	serverBaseURLArg := flag.String("b", "http://localhost:8080", "server base URL")
	fileStoragePathArg := flag.String("f", "/tmp/short-url-db.json", "file storage path")

	databaseDSNArg := flag.String("d", "", "database DSN")



	// делаем разбор командной строки
	flag.Parse()

	// переменная для хранения настроек конфигурации
	var cfg Config

	// парсим переменные среды в структуру
	err = env.Parse(&cfg)
	if err != nil {
		logger.Fatal("error occured while Parse env: " + err.Error())
	}

	// адрес сервера из флага
	serverAddress := *serverAddressArg

	// базовый URL из флага
	serverBaseURL := *serverBaseURLArg

	// базовый URL из флага
	fileStoragePath := *fileStoragePathArg


	// базовый URL из флага
	databaseDSN := *databaseDSNArg


	// Если переменные среды установлены, то берем данные эти данные
	if cfg.ServerAddress != "" {
		serverAddress = cfg.ServerAddress
	}

	// Если переменные среды установлены, то берем данные эти данные
	if cfg.BaseURL != "" {
		serverBaseURL = cfg.BaseURL
	}

	// Если переменные среды установлены, то берем данные эти данные
	if cfg.FileStoragePath != "" {
		fileStoragePath = cfg.FileStoragePath
	}


	// Если переменные среды установлены, то берем данные эти данные
	if cfg.DatabaseDSN != "" {
		databaseDSN = cfg.DatabaseDSN
	}

	// создаем сервер
	srv := internal.NewServer(serverAddress, serverBaseURL, fileStoragePath, databaseDSN, logger)


	// запускаем сервер
	if err := srv.Run(srv.GetServerAddress()); err != nil {
		logger.Fatal("error occured while running http server: " + err.Error())
	}

}
