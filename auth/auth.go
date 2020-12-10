package auth

import (
	"context"
	"errors"
	"files-back/handlers"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	SecretKey  []byte
	TokenTTL   = time.Hour * 24
	userCtxKey = contextKey("username")
)

var (
	Unauthorized = errors.New("unauthorized users")
	BadToken     = errors.New("bad token")
)

type contextKey string

type incomingJSON struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Token struct {
	Token   string `json:"token"`
	Expires string `json:"expires"`
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), " ")
		if len(authHeader) != 2 {
			handlers.StatusUnauthorized(Unauthorized, w)
			return
		}
		username, err := ParseToken(authHeader[1])
		if err != nil {
			handlers.StatusUnauthorized(err, w)
			return
		}
		ctx := context.WithValue(r.Context(), userCtxKey, &username)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func GenerateToken(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = time.Now().Add(TokenTTL).Unix()
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		log.Fatal("Error in Generating key")
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(
		tokenStr,
		func(token *jwt.Token) (interface{}, error) {
			return SecretKey, nil
		},
	)
	if token == nil {
		return "", BadToken
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["username"].(string)
		return username, nil
	} else {
		return "", err
	}
}
