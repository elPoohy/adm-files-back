package auth

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	SecretKey  []byte
	TokenTTL   = time.Hour * 24
	userCtxKey = contextKey("user")
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

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("authorization")
			bearerToken := strings.Split(authHeader, " ")
			fmt.Println(authHeader)
			if len(bearerToken) != 2 {
				next.ServeHTTP(w, r)
				return
			}
			username, err := ParseToken(bearerToken[1])
			if err != nil {
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), userCtxKey, &username)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func ForContext(ctx context.Context) *string {
	raw, _ := ctx.Value(userCtxKey).(*string)
	return raw
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
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["username"].(string)
		return username, nil
	} else {
		return "", err
	}
}
