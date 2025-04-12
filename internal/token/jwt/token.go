package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"net/http"
	"projectgrom/web/internal/config"
	"time"
)

var (
	CreateErr = errors.New("create token failed")
	TokenErr  = errors.New("get token failed")
)

type Token struct {
	Name string
	jwt.StandardClaims
}

// InitToken - функция, создаем структуру.
func InitToken(name string) *Token {
	return &Token{
		Name: name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 5).Unix(),
		},
	}
}

// CreateToken - функция, создает токен.
func CreateToken(name string) (string, error) {
	newToken := InitToken(name)
	key := config.SecretKey
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, newToken).SignedString([]byte(key))
	if err != nil {
		return "", fmt.Errorf("%w", CreateErr)
	}
	return token, nil
}

// ParseToken - фукнция, парсит токен.
func ParseToken(tokenString string) (*Token, error) {
	key := config.SecretKey
	token, err := jwt.ParseWithClaims(tokenString, &Token{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return []byte(key), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Token); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}

func ClearToken(w http.ResponseWriter) http.ResponseWriter {
	http.SetCookie(w, &http.Cookie{
		Name:   "Authorization",
		Value:  "",
		MaxAge: -1,
		Path:   "/api/main",
	})
	return w
}
