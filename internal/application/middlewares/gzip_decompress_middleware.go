package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// compressReader реализует интерфейс io.ReadCloser и позволяет прозрачно для сервера
// декомпрессировать получаемые от клиента данные
type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

func GzipDecompressMiddleware(h http.Handler) http.Handler {
	gzipFn := func(rw http.ResponseWriter, req *http.Request) {

		ow := rw
		// проверяем, что клиент отправил серверу сжатые данные в формате gzip
		contentEncoding := req.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
			cr, err := newCompressReader(req.Body)
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			// меняем тело запроса на новое
			req.Body = cr
			defer cr.Close()
		}

		// передаём управление хендлеру
		h.ServeHTTP(ow, req)

		/*
			log.Printf("Header: Content-Encoding: %v", req.Header.Get("Content-Encoding"))

			if !strings.Contains(req.Header.Get("Content-Encoding"), "gzip") {
				// если gzip не поддерживается, передаём управление дальше без изменений
				h.ServeHTTP(rw, req)
				return
			}

			// создаём *gzip.Reader, который будет читать тело запроса и распаковывать его
			gz, err := gzip.NewReader(req.Body)
			if err != nil {
				// log.Println(err.Error())
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}

			log.Println("test 2")
			// закрываем gzip-читателя
			defer gz.Close()

			// при чтении вернётся распакованный слайс байт
			body, err := io.ReadAll(gz)
			log.Println("test 3")
			// log.Printf("Body: %vr\n", req.Header.Get("Content-Encoding"))
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}

			// после того как распаовали тело запроса,
			// нужно его перезаписать
			req.Body = io.NopCloser(strings.NewReader(string(body)))
			req.ContentLength = int64(len(string(body)))
		*/
	}

	return http.HandlerFunc(gzipFn)
}
