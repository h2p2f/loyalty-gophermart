package httpserver

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"

	j "github.com/h2p2f/loyalty-gophermart/internal/gophermart/utils/jwt"
)

// contextKey is a value for use with context.WithValue.
type contextKey string

// String returns the string representation of the context key.
func (c contextKey) String() string {
	return string(c)
}

// loginContextKey is the context key for the login.
var loginContextKey = contextKey("login")

// LoginFromContext returns the login value stored in ctx
func LoginFromContext(ctx context.Context) (string, bool) {
	login, ok := ctx.Value(loginContextKey).(string)
	return login, ok
}

// JWTAuth is a middleware that checks for a valid JWT token in the Authorization header.
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
		// Get the JWT token from the Authorization header
		tokenHeader := r.Header.Get("Authorization")
		// if token is missing, returns error code 403 Unauthorized
		if tokenHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// The token normally comes in format `Bearer {token-body}`
		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// Get the token part
		tokenPart := splitted[1]
		// Initialize a new instance of `Claims`
		tk := &j.Claims{}
		// Parse the JWT token and claims
		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(j.SECRET), nil
		})

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// Token is invalid, maybe not signed on this server
		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// write login to context
		ctx := context.WithValue(r.Context(), loginContextKey, tk.Login)
		r = r.WithContext(ctx)
		w.Header().Add("Authorization", tokenHeader)
		// next middleware chain
		next.ServeHTTP(w, r)
	})

}

// GzipHanler is middleware for gzip
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
