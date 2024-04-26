package internal

import (
	"net/http"

	"github.com/cfif1982/urlshtr.git/pkg/log"

	"github.com/cfif1982/urlshtr.git/internal/application/handlers"
	"github.com/cfif1982/urlshtr.git/internal/application/middlewares"
	linksInfra "github.com/cfif1982/urlshtr.git/internal/infrastructure/links"

	"github.com/go-chi/chi/v5"
)

// структура сервера. Храним передаваемые параметры при запуске программы
type Server struct {
	serverAddress string
	serverBaseURL string
	logger        *log.Logger
}

// Конструктор Server
func NewServer(addr string, base string, logger *log.Logger) Server {
	return Server{
		serverAddress: addr,
		serverBaseURL: base,
		logger:        logger,
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
	// Создаем репозиторий для работы с БД. Здесь можно изменить БД и выбратьдругую тенологию
	linkRepo := linksInfra.NewLocalRepository()

	// создаем хндлер и передаем ему нужную БД
	handler := handlers.NewHandler(linkRepo, s.serverBaseURL, s.logger)
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
	router.Post(`/`, middlewares.LogMiddleware(s.logger, http.HandlerFunc(handler.AddLink)))
	router.Post(`/api/shorten`, middlewares.LogMiddleware(s.logger, http.HandlerFunc(handler.PostAddLink)))

	return router
}
