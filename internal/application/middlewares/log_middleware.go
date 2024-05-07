package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cfif1982/urlshtr.git/pkg/log"
)

type (
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

// переопределяем функцию Write дл получения размера записанных данных
func (r *loggingResponseWriter) Write(b []byte) (int, error) {

	// вызываем оригинальную функцию Write
	size, err := r.ResponseWriter.Write(b)

	// схраняем размер записанных данных
	r.resData.size = size

	return size, err
}

// переопределяем функцию WriteHeader для получения кода ответа
func (r *loggingResponseWriter) WriteHeader(statusCode int) {

	// вызываем оригинальную функцию WriteHeader
	r.ResponseWriter.WriteHeader(statusCode)

	// сохраняем код ответа
	r.resData.status = statusCode
}

// middleware для логирования хэндлеров
func LogMiddleware(logger *log.Logger, h http.Handler) http.HandlerFunc {

	logFn := func(rw http.ResponseWriter, req *http.Request) {

		// запоминаем время начала обработки запроса
		start := time.Now()

		// создаем структуру для хранения нужных данных
		rd := responseData{
			status: 0,
			size:   0,
		}

		// создаем переопределяемую структуру ResponseWriter
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
		logger.Info(
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
