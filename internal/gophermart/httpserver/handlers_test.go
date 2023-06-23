package httpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/h2p2f/loyalty-gophermart/internal/gophermart/httpserver/mocks"
	"github.com/h2p2f/loyalty-gophermart/internal/gophermart/models"
	"github.com/h2p2f/loyalty-gophermart/internal/gophermart/utils/luhn"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGopherMartHandler_AddOrder(t *testing.T) {

	tests := []struct {
		name   string
		user   string
		code   int
		order  string
		exists bool
		owner  string
	}{
		{
			name:   "AddOrder PositiveTest 1",
			user:   "test",
			code:   http.StatusAccepted,
			order:  "7376351891",
			exists: false,
		},
		{
			name:   "AddOrder PositiveTest 2",
			user:   "test",
			code:   http.StatusAccepted,
			order:  "6512136380",
			exists: false,
		},
		{
			name:   "AddOrder Wrong Order Number",
			user:   "test",
			code:   http.StatusUnprocessableEntity,
			order:  "1234567890",
			exists: false,
		},
		{
			name:   "AddOrder if Order Exists and Owner is the Same",
			user:   "test",
			code:   http.StatusOK,
			order:  "6512136380",
			exists: true,
			owner:  "test",
		},
		{
			name:   "AddOrder if Order Exists and Owner is Different",
			user:   "test",
			code:   http.StatusConflict,
			order:  "6512136380",
			exists: true,
			owner:  "test2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWriter := httptest.NewRecorder()
			body := bytes.NewBufferString(tt.order)
			mockRequest := httptest.NewRequest("POST", "/orders", body)
			var mockContext context.Context
			if tt.user != "" {
				mockContext = context.WithValue(context.Background(), loginContextKey, tt.user)
			} else {
				mockContext = context.Background()
			}
			mockDB := mocks.NewDataBaser(t)
			if tt.user != "" && luhn.Validate(tt.order) && !tt.exists {
				mockDB.On("CheckUniqueOrder", mock.Anything, tt.order).Return("", false)
				mockDB.On("NewOrder", mock.Anything, tt.user, mock.Anything).Return(nil)
			} else if tt.user != "" && !luhn.Validate(tt.order) {
				mockDB.On("CheckUniqueOrder", mock.Anything, tt.order).Return("", false)
			}
			if tt.exists {
				mockDB.On("CheckUniqueOrder", mock.Anything, tt.order).Return(tt.owner, true)
			}
			handler := &GopherMartHandler{
				db:     mockDB,
				logger: zap.NewNop(),
			}
			handler.AddOrder(mockWriter, mockRequest.WithContext(mockContext))
			assert.Equal(t, tt.code, mockWriter.Code)

		})
	}
}

func TestGopherMartHandler_Balance(t *testing.T) {
	tests := []struct {
		name        string
		balance     float64
		withdraws   float64
		user        string
		code        int
		contentType string
	}{
		{
			name:        "GetBalance PositiveTest 1",
			balance:     100.1,
			withdraws:   50.2,
			user:        "test_login",
			code:        http.StatusOK,
			contentType: "application/json",
		},
		{
			name:        "GetBalance PositiveTest 2",
			balance:     200.1,
			withdraws:   500000.2,
			user:        "test_login",
			code:        http.StatusOK,
			contentType: "application/json",
		},
		{
			name:        "GetBalance NegativeTest 1",
			balance:     0,
			withdraws:   0,
			user:        "",
			code:        http.StatusInternalServerError,
			contentType: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockWriter := httptest.NewRecorder()
			mockRequest := httptest.NewRequest("GET", "/balance", nil)
			var mockContext context.Context
			if tt.user != "" {
				mockContext = context.WithValue(context.Background(), loginContextKey, tt.user)
			} else {
				mockContext = context.Background()
			}
			mockDB := mocks.NewDataBaser(t)

			if tt.user != "" {
				mockDB.On("GetBalance", mock.Anything, tt.user).Return(tt.balance, nil)
				mockDB.On("GetSumOfAllWithdraws", mock.Anything, tt.user).Return(tt.withdraws)
			}
			handler := &GopherMartHandler{
				db:     mockDB,
				logger: zap.NewNop(),
			}
			handler.GetBalance(mockWriter, mockRequest.WithContext(mockContext))
			assert.Equal(t, tt.code, mockWriter.Code)
			assert.Equal(t, tt.contentType, mockWriter.Header().Get("Content-Type"))
			if tt.code == 200 {
				expectedAccount := models.Account{Balance: tt.balance, Withdraws: tt.withdraws}
				expectedResponse, _ := json.Marshal(expectedAccount)
				assert.Equal(t, expectedResponse, mockWriter.Body.Bytes())
			}

		})
	}
}

func TestGopherMartHandler_Login(t *testing.T) {

	tests := []struct {
		name              string
		user              string
		password          string
		encryptedPassword string
		wrongUser         bool
		code              int
	}{
		{
			name:              "Login PositiveTest 1",
			user:              "test_login",
			password:          "12345678",
			encryptedPassword: "$2a$10$4s6.ghWw25/q2fxwLNh/N.UVMDNTK/GhQNR9P2JZALP.bX97ttwOe",
			code:              http.StatusOK,
			wrongUser:         false,
		},
		{
			name:              "Login NegativeTest 1 - wrong password",
			user:              "test_login",
			password:          "123456789",
			encryptedPassword: "$2a$10$4s6.ghWw25/q2fxwLNh/N.UVMDNTK/GhQNR9P2JZALP.bX97ttwOe",
			code:              http.StatusUnauthorized,
			wrongUser:         false,
		},
		{
			name:              "Login NegativeTest 2 - wrong user",
			user:              "test_login",
			password:          "12345678",
			encryptedPassword: "$2a$10$4s6.ghWw25/q2fxwLNh/N.UVMDNTK/GhQNR9P2JZALP.bX97ttwOe",
			code:              http.StatusUnauthorized,
			wrongUser:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			u := []byte(`
			{
				"login": "` + tt.user + `",
				"password": "` + tt.password + `"
			}
		`)
			body := bytes.NewBuffer(u)
			mockWriter := httptest.NewRecorder()
			mockRequest := httptest.NewRequest("POST", "/login", body)
			mockDB := mocks.NewDataBaser(t)
			if !tt.wrongUser {
				mockDB.On("FindPassByLogin", mock.Anything, tt.user).Return(tt.encryptedPassword, nil)
			} else {
				mockDB.On("FindPassByLogin", mock.Anything, tt.user).Return("", errors.New("user not found"))
			}
			handler := &GopherMartHandler{
				db:     mockDB,
				logger: zap.NewNop(),
			}
			handler.Login(mockWriter, mockRequest)
			assert.Equal(t, tt.code, mockWriter.Code)
		})
	}
}

func TestGopherMartHandler_Orders(t *testing.T) {

	tests := []struct {
		name        string
		user        string
		code        int
		ordersDasta []models.Order
	}{
		{
			name: "GetOrders PositiveTest",
			user: "test_login",
			code: http.StatusOK,
			ordersDasta: []models.Order{
				{
					Number:      "12345678903",
					Status:      models.NEW,
					Accrual:     500.12,
					TimeCreated: time.Now(),
				},
				{
					Number:      "12345678904",
					Status:      models.PROCESSED,
					Accrual:     700.78,
					TimeCreated: time.Now(),
				},
			},
		},
		{
			name:        "GetOrders not authorized",
			user:        "",
			code:        http.StatusInternalServerError,
			ordersDasta: []models.Order{},
		},
		{
			name:        "GetOrders not found",
			user:        "test_login",
			code:        http.StatusInternalServerError,
			ordersDasta: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWriter := httptest.NewRecorder()
			mockRequest := httptest.NewRequest("GET", "/orders", nil)
			var mockContext context.Context
			if tt.user != "" {
				mockContext = context.WithValue(context.Background(), loginContextKey, tt.user)
			} else {
				mockContext = context.Background()
			}
			mockDB := mocks.NewDataBaser(t)
			mockData, err := json.Marshal(tt.ordersDasta)
			if err != nil {
				fmt.Println(err)
			}
			if tt.user != "" && tt.ordersDasta != nil {
				mockDB.On("GetOrdersByUser", mock.Anything, tt.user).Return(mockData, nil)
			} else if tt.user != "" && tt.ordersDasta == nil {
				mockDB.On("GetOrdersByUser", mock.Anything, tt.user).Return(nil, errors.New("not found"))
			}

			handler := &GopherMartHandler{
				db:     mockDB,
				logger: zap.NewNop(),
			}
			handler.GetOrders(mockWriter, mockRequest.WithContext(mockContext))
			assert.Equal(t, tt.code, mockWriter.Code)
			if tt.code == http.StatusOK {
				assert.Equal(t, mockData, mockWriter.Body.Bytes())
			}
		})
	}
}

func TestGopherMartHandler_Register(t *testing.T) {

	tests := []struct {
		name     string
		user     string
		password string
		ifExists bool
		code     int
	}{
		{
			name:     "Register PositiveTest",
			user:     "test_login",
			password: "123456789",
			ifExists: false,
			code:     http.StatusOK,
		},
		{
			name:     "Register NegativeTest - user exists",
			user:     "test_login",
			password: "123456789",
			ifExists: true,
			code:     http.StatusConflict,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pgErr pgconn.PgError
			pgErr.Message = "duplicate login key"
			pgErr.Severity = "ERROR"
			pgErr.Code = "23505"
			err := &pgErr
			mockDB := mocks.NewDataBaser(t)
			u := []byte(`
			{
				"login": "` + tt.user + `",
				"password": "` + tt.password + `"
			}
		`)
			body := bytes.NewBuffer(u)
			mockWriter := httptest.NewRecorder()
			mockRequest := httptest.NewRequest("POST", "/register", body)
			if tt.ifExists {
				mockDB.On("NewUser", mock.Anything, mock.Anything).Return(err)
			} else {
				mockDB.On("NewUser", mock.Anything, mock.Anything).Return(nil)
			}
			handler := &GopherMartHandler{
				db:     mockDB,
				logger: zap.NewNop(),
			}
			handler.Register(mockWriter, mockRequest)
			assert.Equal(t, tt.code, mockWriter.Code)
			if tt.code == http.StatusOK {
				token := mockWriter.Header().Get("Authorization")
				split := strings.Split(token, " ")
				assert.Equal(t, "Bearer", split[0])
			}
		})
	}
}

func TestGopherMartHandler_Withdraw(t *testing.T) {

	tests := []struct {
		name    string
		user    string
		order   string
		sum     float64
		balance float64
		code    int
	}{
		{
			name:    "DoWithdraw PositiveTest",
			user:    "test_login",
			order:   "2377225624",
			sum:     500.12,
			balance: 500.6,
			code:    http.StatusOK,
		},
		{
			name:    "DoWithdraw NegativeTest - not enough money",
			user:    "test_login",
			order:   "2377225624",
			sum:     500.12,
			balance: 0.6,
			code:    http.StatusPaymentRequired,
		},
		{
			name:    "DoWithdraw NegativeTest - not authorized",
			user:    "",
			order:   "2377225624",
			sum:     500.12,
			balance: 0.6,
			code:    http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := mocks.NewDataBaser(t)
			mockWriter := httptest.NewRecorder()
			body := bytes.NewBuffer([]byte(`
			{	
				"order": "` + tt.order + `",
				"sum": ` + fmt.Sprintf("%f", tt.sum) + `
			}	
		`))
			mockRequest := httptest.NewRequest("POST", "/withdraw", body)
			var mockContext context.Context
			if tt.user != "" {
				mockContext = context.WithValue(context.Background(), loginContextKey, tt.user)
			} else {
				mockContext = context.Background()
			}
			if tt.user != "" && tt.balance > tt.sum {
				mockDB.On("GetBalance", mock.Anything, tt.user).Return(tt.balance, nil)
				mockDB.On("NewOrder", mock.Anything, tt.user, mock.Anything).Return(nil)
				mockDB.On("NewWithdraw", mock.Anything, tt.user, mock.Anything).Return(nil)
			} else if tt.user != "" && tt.balance < tt.sum {
				mockDB.On("GetBalance", mock.Anything, tt.user).Return(tt.balance, nil)
			}
			handler := &GopherMartHandler{
				db:     mockDB,
				logger: zap.NewNop(),
			}
			handler.DoWithdraw(mockWriter, mockRequest.WithContext(mockContext))
			assert.Equal(t, tt.code, mockWriter.Code)
		})
	}
}

func TestGopherMartHandler_Withdrawals(t *testing.T) {

	tests := []struct {
		name string
		user string
		code int
	}{
		{
			name: "GetWithdrawals PositiveTest",
			user: "test_login",
			code: http.StatusOK,
		},
		{
			name: "GetWithdrawals NegativeTest - not authorized",
			user: "",
			code: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := mocks.NewDataBaser(t)
			var mockContext context.Context
			if tt.user != "" {
				mockContext = context.WithValue(context.Background(), loginContextKey, tt.user)
			} else {
				mockContext = context.Background()
			}
			if tt.user != "" {
				someString := "some string"
				someStringBytes := []byte(someString)
				mockDB.On("GetAllWithdraws", mock.Anything, tt.user).Return(someStringBytes, nil)
			}
			mockWriter := httptest.NewRecorder()
			mockRequest := httptest.NewRequest("GET", "/withdrawals", nil)

			handler := &GopherMartHandler{
				db:     mockDB,
				logger: zap.NewNop(),
			}
			handler.GetWithdrawals(mockWriter, mockRequest.WithContext(mockContext))
			assert.Equal(t, tt.code, mockWriter.Code)
		})
	}
}
