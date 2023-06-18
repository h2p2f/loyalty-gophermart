package httpserver

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"

	j "github.com/h2p2f/loyalty-gophermart/internal/gophermart/utils/jwt"
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

var loginContextKey = contextKey("login")

func LoginFromContext(ctx context.Context) (string, bool) {
	login, ok := ctx.Value(loginContextKey).(string)
	return login, ok
}

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
		ctx := context.WithValue(r.Context(), loginContextKey, tk.Login)
		r = r.WithContext(ctx)
		w.Header().Add("Authorization", tokenHeader)
		next.ServeHTTP(w, r)
	})

}

func GzipHanler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		acceptEncoding := r.Header.Get("Accept-Encoding")
		contentEncoding := r.Header.Get("Content-Encoding")
		supportGzip := strings.Contains(acceptEncoding, "gzip")
		sendGzip := strings.Contains(contentEncoding, "gzip")
		// this section for ordinary request
		if !supportGzip && !sendGzip {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
			return
		}
		// this section for request with accept-encoding: gzip
		if supportGzip && !sendGzip {
			originWriter := w
			compressedWriter := NewCompressWriter(w)

			originWriter = compressedWriter
			originWriter.Header().Set("Content-Encoding", "gzip")
			defer compressedWriter.Close()

			next.ServeHTTP(originWriter, r)
		}
		// this section for request with content-encoding: gzip
		if sendGzip {
			originWriter := w
			compressedWriter := NewCompressWriter(w)
			originWriter = compressedWriter
			originWriter.Header().Set("Content-Encoding", "gzip")
			defer compressedWriter.Close()

			compressedReader, err := NewCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			r.Body = compressedReader
			defer compressedReader.Close()

			next.ServeHTTP(originWriter, r)

		}
	})

}
