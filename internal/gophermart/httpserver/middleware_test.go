package httpserver

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	j "github.com/h2p2f/loyalty-gophermart/internal/gophermart/utils/jwt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJWTAuth(t *testing.T) {

	tests := []struct {
		name     string
		user     string
		url      string
		code     int
		notValid bool
	}{
		{
			name:     "Positive test",
			user:     "test_login",
			url:      "/api/user/orders",
			code:     http.StatusOK,
			notValid: false,
		},
		{
			name:     "Negative test - unauthorized",
			user:     "test_login",
			url:      "/api/user/orders",
			code:     http.StatusUnauthorized,
			notValid: false,
		},
		{
			name:     "Negative test - invalid token",
			user:     "test_login",
			url:      "/api/user/orders",
			code:     http.StatusUnauthorized,
			notValid: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWriter := httptest.NewRecorder()
			mockRequest, _ := http.NewRequest(http.MethodGet, tt.url, nil)
			ctx := context.WithValue(context.Background(), secretContextKey, "somesecretkey")
			mockRequest = mockRequest.WithContext(ctx)
			if tt.code == 200 {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, &j.Claims{Login: tt.user})
				tokenString, _ := token.SignedString([]byte("somesecretkey"))
				mockRequest.Header.Set("Authorization", "Bearer "+tokenString)
			}
			if tt.notValid {
				mockRequest.Header.Set("Authorization", "Bearer invalid_token")
			}
			mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "test_login", r.Context().Value(loginContextKey))
				w.WriteHeader(http.StatusOK)
			})
			handler := JWTAuth(mockHandler)
			handler.ServeHTTP(mockWriter, mockRequest)
			assert.Equal(t, tt.code, mockWriter.Code)
		})
	}
}
