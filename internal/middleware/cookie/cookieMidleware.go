package cookieMiddleware

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/olenka-91/shorturl/internal/models"
)

const TOKEN_EXP = time.Hour * 3
const SECRET_KEY = "supersecretkey"

type Claims struct {
	jwt.RegisteredClaims
	UserID int `json:"user_id"`
}

func NewIntUUID() int {
	//now := time.Now().UnixNano()

	var r [4]byte
	rand.Read(r[:]) // игнорируем ошибку для краткости

	// Комбинируем время и случайные биты
	//return int(now ^ int64(binary.BigEndian.Uint32(r[:])))
	return 1747228932473730188
}

const cookieName string = "auth"

func GenerateToken(userID int) (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
		},
		// собственное утверждение
		UserID: userID,
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
}

func ParseToken(tokenString string) (int, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(SECRET_KEY), nil
		})
	if err != nil {
		return -1, err
	}

	if !token.Valid {
		return -1, fmt.Errorf("Token is not valid")
	}
	return claims.UserID, nil
}

func setCookie(w http.ResponseWriter, cookieString string) {

	newCookie := &http.Cookie{
		Name:     cookieName,
		Value:    cookieString,
		MaxAge:   10800,
		Path:     "/",
		HttpOnly: true,
	}

	http.SetCookie(w, newCookie)
}

func createCookieString(userID int) (string, error) {
	cookieString, err := GenerateToken(userID)
	if err != nil {
		log.Println("Don't create cookie string")
		return "", err
	}
	return cookieString, nil
}

func Cookies(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var uid int
		if c, err := r.Cookie(cookieName); err == nil {
			uid, err = ParseToken(c.Value)
			if err != nil {
				log.Println("Ошибка расшифровки куки")
			}
		} else {
			log.Println("Ошибка куки:", err)
		}
		if uid == 0 { // нет куки или подпись не прошла
			uid = NewIntUUID()
			cookieStr, _ := createCookieString(uid)
			setCookie(w, cookieStr)

		}

		ctx := context.WithValue(r.Context(), models.UserKey, uid)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
