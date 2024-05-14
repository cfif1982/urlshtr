package internal

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cfif1982/urlshtr.git/pkg/log"

	"github.com/cfif1982/urlshtr.git/internal/application/handlers"
	"github.com/cfif1982/urlshtr.git/internal/application/middlewares"
	linksInfra "github.com/cfif1982/urlshtr.git/internal/infrastructure/links"

	"github.com/go-chi/chi/v5"
)

// структура сервера. Храним передаваемые параметры при запуске программы
type Server struct {
	serverAddress   string
	serverBaseURL   string
	FileStoragePath string
	databaseDSN     string
	logger          *log.Logger
}

// Конструктор Server
func NewServer(addr, base, storage, dsn string, logger *log.Logger) Server {
	return Server{
		serverAddress:   addr,
		serverBaseURL:   base,
		FileStoragePath: storage,
		databaseDSN:     dsn,
		logger:          logger,
	}
}

// получить адресс сервера
func (s *Server) GetServerAddress() string {
	return s.serverAddress
}

// получить базовый URL
func (s *Server) GetServerBaseURL() string {
	return s.serverBaseURL
}

// запуск сервера
func (s *Server) Run(serverAddr string) error {

	// Dependency Injection
	//********************************************************
	// Создаем репозиторий для работы с БД. Здесь можно изменить БД и выбрать другую тенологию
	// Если указан файл , то база данных х файлов

	var (
		linkRepo handlers.RepositoryInterface
		handler  *handlers.Handler
		err      error
	)

	// для тестирования СУБД на локальном компе
	s.databaseDSN = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", `localhost`, `postgres`, `123`, `videos`)

	// создаю контекст для подключения БД
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// если указан адрес СУБД
	if s.databaseDSN != "" {
		linkRepo, err = linksInfra.NewPostgresRepository(ctx, s.databaseDSN, s.logger)

		if err != nil {
			s.logger.Fatal("can't initialize postgres DB: " + err.Error())
		} else {
			s.logger.Info("postgres DB initialized")
		}
	} else {
		// Если не указан файл как БД, то создаем репозиторий в памяти
		if s.FileStoragePath == "" {
			linkRepo = linksInfra.NewLocalRepository()
			s.logger.Info("Local Repository initialized")
		} else {
			// если указан файл как БД, то инициализируем файловый репозиторий
			linkRepo, err = linksInfra.NewFileRepository(s.FileStoragePath)

			if err != nil {
				s.logger.Fatal("can't initialize storage file: " + err.Error())
			} else {
				s.logger.Info("file storage initialized")
			}
		}
	}

	// создаем хндлер и передаем ему нужную БД
	handler = handlers.NewHandler(linkRepo, s.serverBaseURL, s.logger)

	//********************************************************

	// инициализируем роутер
	routerChi := s.InitRoutes(handler)

	s.logger.Info("Starting server", "addr", serverAddr)

	// запуск сервера на нужно адресе и с нужным роутером
	return http.ListenAndServe(serverAddr, routerChi)
}

// инициализируем роутер CHI
func (s *Server) InitRoutes(handler *handlers.Handler) *chi.Mux {

	// создаем роутер
	router := chi.NewRouter()

	router.Use(middlewares.GzipCompressMiddleware)
	router.Use(middlewares.GzipDecompressMiddleware)

	// назначаем хэндлеры для обработки запросов пользователя
	router.Get(`/{key}`, middlewares.LogMiddleware(s.logger, http.HandlerFunc(handler.GetLinkByKey)))
	router.Get(`/ping`, middlewares.LogMiddleware(s.logger, http.HandlerFunc(handler.Ping)))
	router.Post(`/`, middlewares.LogMiddleware(s.logger, http.HandlerFunc(handler.AddLink)))
	router.Post(`/api/shorten`, middlewares.LogMiddleware(s.logger, http.HandlerFunc(handler.PostAddLink)))
	router.Post(`/api/shorten/batch`, middlewares.LogMiddleware(s.logger, http.HandlerFunc(handler.PostAddBatchLink)))

	return router
}
