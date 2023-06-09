package httpserver

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"

	j "github.com/h2p2f/loyalty-gophermart/internal/gophermart/utils/jwt"
)

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		notAuth := []string{"/api/user/register", "/api/user/login"}
		requestPath := r.URL.Path
		for _, value := range notAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}
		tokenHeader := r.Header.Get("Authorization")

		if tokenHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tokenPart := splitted[1]
		tk := &j.Claims{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(j.SECRET), nil
		})

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		//code for statictests
		type keyType string
		var key keyType = "login"
		ctx := context.WithValue(r.Context(), key, tk.Login)
		r = r.WithContext(ctx)
		w.Header().Add("Authorization", tokenHeader)
		next.ServeHTTP(w, r)
	})

}
