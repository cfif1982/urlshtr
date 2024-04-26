package middlewares

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
)

func GzipDecompressMiddleware(h http.Handler) http.Handler {
	gzipFn := func(rw http.ResponseWriter, req *http.Request) {

		log.Printf("Header: Content-Encoding: %v", req.Header.Get("Content-Encoding"))

		if !strings.Contains(req.Header.Get("Content-Encoding"), "gzip") {
			// если gzip не поддерживается, передаём управление дальше без изменений
			h.ServeHTTP(rw, req)
			return
		}

		// создаём *gzip.Reader, который будет читать тело запроса и распаковывать его
		gz, err := gzip.NewReader(req.Body)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		// закрываем gzip-читателя
		defer gz.Close()

		// при чтении вернётся распакованный слайс байт
		body, err := io.ReadAll(gz)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		// после того как распаовали тело запроса,
		// нужно его перезаписать
		req.Body = io.NopCloser(strings.NewReader(string(body)))
		req.ContentLength = int64(len(string(body)))

	}

	return http.HandlerFunc(gzipFn)
}
