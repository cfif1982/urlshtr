package internal

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cfif1982/urlshtr.git/pkg/log"

	"github.com/cfif1982/urlshtr.git/internal/application/handlers"
	linksInfra "github.com/cfif1982/urlshtr.git/internal/infrastructure/links"

	"github.com/go-chi/chi/v5"
)

type (
	// структура сервера. Храним передаваемые параметры при запуске программы
	Server struct {
		serverAddress string
		serverBaseURL string
		logger        *log.Logger
	}

	// структура для хранения данных о параметрах ответа сервера
	responseData struct {
		status int
		size   int
	}

	// своя реализация интерфейса ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter
		resData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {

	size, err := r.ResponseWriter.Write(b)
	r.resData.size = size

	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {

	r.ResponseWriter.WriteHeader(statusCode)

	r.resData.status = statusCode
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

	// назначаем хэндлеры для обработки запросов пользователя
	router.Get(`/{key}`, s.middlewareLogging(http.HandlerFunc(handler.GetLinkByKey)))
	router.Post(`/`, s.middlewareLogging(http.HandlerFunc(handler.AddLink)))

	return router
}

func (s *Server) middlewareLogging(h http.Handler) http.HandlerFunc {

	logFn := func(rw http.ResponseWriter, req *http.Request) {

		// апоминаем время начала обработки запроса
		start := time.Now()

		rd := responseData{
			status: 0,
			size:   0,
		}

		logRW := loggingResponseWriter{
			ResponseWriter: rw,
			resData:        &rd,
		}

		// нужные переменные для вывода в логе
		uri := req.RequestURI
		method := req.Method

		// выполняем оригинальный запрос
		// вот тут не понял((
		// почему  аргумент logRW нужно передавать по ссылке?
		// h.ServeHTTP(logRW, req) - выдает ошибку
		h.ServeHTTP(&logRW, req)

		// вычисляем время выполнения запроса
		duration := time.Since(start)

		// выводим лог
		s.logger.Info(
			"request info:",
			"uri", uri,
			"method", method,
			"status", fmt.Sprint(rd.status),
			"duration", duration.String(),
			"size", fmt.Sprint(rd.size),
		)
	}

	return http.HandlerFunc(logFn)
}
