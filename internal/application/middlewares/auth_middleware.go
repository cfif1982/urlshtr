package middlewares

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Claims — структура утверждений, которая включает стандартные утверждения и
// одно пользовательское UserID
type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

const TOKEN_EXP = time.Hour * 3
const COOKIE_NAME = "accessToken"

func AuthMiddleware(h http.Handler) http.Handler {
	authFn := func(rw http.ResponseWriter, req *http.Request) {

		// получаем токен из куки
		tokenFromCookie, err := req.Cookie(COOKIE_NAME)

		// если такой куки нет, то создаем куку и выдаем сообщение об этом и сттус http.StatusUnauthorized
		if err != nil {
			// строим строку токена для куки
			token, _ := buildJWTString()

			// создаем куку в http
			cookie := http.Cookie{}
			cookie.Name = COOKIE_NAME
			cookie.Value = token
			cookie.Expires = time.Now().Add(TOKEN_EXP)
			cookie.Path = "/"

			// устанавливаем созданную куку в http
			http.SetCookie(rw, &cookie)

			// выводим сообщение об ошибке
			http.Error(rw, "token missed", http.StatusUnauthorized)

			return
		}

		// получаем user id из токена
		userId := getUserIdFromToken(tokenFromCookie.Value)

		// если в токенен нет узера, то сообщаем об этом
		if userId == -1 {
			// выводим сообщение об ошибке и статус http.StatusNoContent
			http.Error(rw, "token missed", http.StatusNoContent)

			return
		}

		// создаю контекст для сохранения userID
		ctx := context.WithValue(req.Context(), "userID", userId)

		// обрабатываем сам запрос
		h.ServeHTTP(rw, req.WithContext(ctx))
	}

	return http.HandlerFunc(authFn)
}

// получаем user id из токена
func getUserIdFromToken(tokenString string) int {
	claims := &Claims{}

	// получаем ключ для генерации токена
	key := getKeyForTokenGeneration()

	// парсим из строки токена tokenString в структуру claims
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return key, nil
		})

	if err != nil {
		return -1
	}

	if !token.Valid {
		return -1
	}

	return claims.UserID
}

// строим строку для токена
func buildJWTString() (string, error) {

	// генерируем user id
	userId, err := createUserID()
	if err != nil {
		return "", err
	}

	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
		},
		UserID: userId,
	})

	// получаем ключ для генерации токена
	key := getKeyForTokenGeneration()

	// создаём строку токена
	strToken, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return strToken, nil
}

// генерируем user id
func createUserID() (int, error) {

	// генерируем случайную последовательность из 6 байт
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		return -1, err
	}

	// преобразуем байты в число
	userID := int(binary.BigEndian.Uint32(b))

	return userID, nil
}

// получаем ключ для генерции токена
func getKeyForTokenGeneration() []byte {

	SECRET_KEY := "supersecretkey"

	return []byte(SECRET_KEY)
}
