package internal

import (
	"net/http"

	"github.com/cfif1982/urlshtr.git/internal/application/handlers"
	linksInfra "github.com/cfif1982/urlshtr.git/internal/infrastructure/links"

	"github.com/go-chi/chi/v5"
)

// структура сервера. Храним передаваемые параметры при запуске программы
type Server struct {
	serverAddress string
	serverBaseURL string
}

// получить адресс сервера
func (s *Server) ServerAddress() string {
	return s.serverAddress
}

// получить базовый URL
func (s *Server) ServerBaseURL() string {
	return s.serverBaseURL
}

// установить адресс сервера
func (s *Server) SetServerAddress(addr string) {
	s.serverAddress = addr
}

// установить базовый URL
func (s *Server) SetServerBaseURL(base string) {
	s.serverBaseURL = base
}

// запуск сервера
func (s *Server) Run(serverAddr string) error {

	// Dependency Injection
	//********************************************************
	// Создаем репозиторий для работы с БД. Здесь можно изменить БД и выбратьдругую тенологию
	linkRepo := linksInfra.NewLocalRepository()

	// создаем хндлер и передаем ему нужную БД
	handler := handlers.NewHandler(linkRepo, s.serverBaseURL)
	//********************************************************

	// инициализируем роутер
	// routerChi := s.InitRoutes(handler)
	//**************************************************
	// создаем роутер
	routerChi := chi.NewRouter()

	// назначаем хэндлеры для обработки запросов пользователя
	routerChi.Get(`/{key}`, handler.GetLinkByKey)
	routerChi.Post(`/`, handler.AddLink)
	//**************************************************

	// запуск сервера на нужно адресе и с нужным роутером
	return http.ListenAndServe(serverAddr, routerChi)
}

//methods on the same type should have the same receiver name (seen 1x "h", 5x "s")

// инициализируем роутер CHI
func (h *Server) Ir(handler *handlers.Handler) *chi.Mux {

	// создаем роутер
	router := chi.NewRouter()

	// назначаем хэндлеры для обработки запросов пользователя
	router.Get(`/{key}`, handler.GetLinkByKey)
	router.Post(`/`, handler.AddLink)

	return router
	// return nil
	//"go.lintFlags": ["--fast"]

}
