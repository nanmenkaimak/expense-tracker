package main

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

var jwtKey = []byte("my_secret_key")

func parseToken(tokenString string) (jwt.MapClaims, error) {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("wrong token")
		}
		return jwtKey, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("token is not valid")
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			http.Error(w, errors.New("no auth header").Error(), http.StatusUnauthorized)
			return
		}
		headerInSlice := strings.Split(header, " ")
		if len(headerInSlice) != 2 {
			http.Error(w, errors.New("wrong header").Error(), http.StatusBadRequest)
			return
		}
		claims, err := parseToken(headerInSlice[1])
		if err != nil {
			http.Error(w, errors.Wrap(err, "parse token").Error(), http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), "id", claims["id"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
