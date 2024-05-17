package middlewares

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"net/http"
	"time"

	"github.com/cfif1982/urlshtr.git/internal/application/handlers"

	"github.com/golang-jwt/jwt/v4"
)

// Claims — структура утверждений, которая включает стандартные утверждения и
// одно пользовательское UserID
type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

const TokenEXP = time.Hour * 3
const CookieName = "accessToken"

func AuthMiddleware(h http.Handler) http.Handler {
	authFn := func(rw http.ResponseWriter, req *http.Request) {

		// получаем токен из куки
		tokenFromCookie, err := req.Cookie(CookieName)

		// если такой куки нет, то создаем куку
		if err != nil {
			cookie := createCookie()

			// устанавливаем созданную куку в http
			http.SetCookie(rw, cookie)

			// обрабатываем запрос
			h.ServeHTTP(rw, req)

			return
		}

		// получаем user id из токена
		userID := getUserIDFromToken(tokenFromCookie.Value)

		// если в токенен нет узера, то pfyjdj cjplftv rere
		if userID == -1 {
			cookie := createCookie()

			// устанавливаем созданную куку в http
			http.SetCookie(rw, cookie)

			// обрабатываем запрос
			h.ServeHTTP(rw, req)
		}

		// создаю контекст для сохранения userID
		ctx := context.WithValue(req.Context(), handlers.KeyUserID, userID)

		// обрабатываем сам запрос
		h.ServeHTTP(rw, req.WithContext(ctx))
	}

	return http.HandlerFunc(authFn)
}

// создаем куку
func createCookie() *http.Cookie {

	// строим строку токена для куки
	token, _ := buildJWTString()

	// создаем куку в http
	cookie := http.Cookie{}
	cookie.Name = CookieName
	cookie.Value = token
	cookie.Expires = time.Now().Add(TokenEXP)
	cookie.Path = "/"

	return &cookie
}

// получаем user id из токена
func getUserIDFromToken(tokenString string) int {
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
	userID, err := createUserID()
	if err != nil {
		return "", err
	}

	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenEXP)),
		},
		UserID: userID,
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

	SecretKEY := "supersecretkey"

	return []byte(SecretKEY)
}
