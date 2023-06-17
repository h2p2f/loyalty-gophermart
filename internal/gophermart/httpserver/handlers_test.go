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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"reflect"
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
			name:        "Balance PositiveTest 1",
			balance:     100.1,
			withdraws:   50.2,
			user:        "test_login",
			code:        200,
			contentType: "application/json",
		},
		{
			name:        "Balance PositiveTest 2",
			balance:     200.1,
			withdraws:   500000.2,
			user:        "test_login",
			code:        200,
			contentType: "application/json",
		},
		{
			name:        "Balance NegativeTest 1",
			balance:     0,
			withdraws:   0,
			user:        "",
			code:        500,
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
			handler.Balance(mockWriter, mockRequest.WithContext(mockContext))
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
			code:              200,
			wrongUser:         false,
		},
		{
			name:              "Login NegativeTest 1 - wrong password",
			user:              "test_login",
			password:          "123456789",
			encryptedPassword: "$2a$10$4s6.ghWw25/q2fxwLNh/N.UVMDNTK/GhQNR9P2JZALP.bX97ttwOe",
			code:              401,
			wrongUser:         false,
		},
		{
			name:              "Login NegativeTest 2 - wrong user",
			user:              "test_login",
			password:          "12345678",
			encryptedPassword: "$2a$10$4s6.ghWw25/q2fxwLNh/N.UVMDNTK/GhQNR9P2JZALP.bX97ttwOe",
			code:              401,
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
			name: "Orders PositiveTest",
			user: "test_login",
			code: 200,
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
			name:        "Orders not authorized",
			user:        "",
			code:        500,
			ordersDasta: []models.Order{},
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
			if tt.user != "" {
				mockDB.On("GetOrdersByUser", mock.Anything, tt.user).Return(mockData, nil)
			}

			handler := &GopherMartHandler{
				db:     mockDB,
				logger: zap.NewNop(),
			}
			handler.Orders(mockWriter, mockRequest.WithContext(mockContext))
			assert.Equal(t, tt.code, mockWriter.Code)
			if tt.code == 200 {
				assert.Equal(t, mockData, mockWriter.Body.Bytes())
			}
		})
	}
}

func TestGopherMartHandler_Register(t *testing.T) {
	type fields struct {
		db     DataBaser
		logger *zap.Logger
	}
	type args struct {
		writer  http.ResponseWriter
		request *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &GopherMartHandler{
				db:     tt.fields.db,
				logger: tt.fields.logger,
			}
			h.Register(tt.args.writer, tt.args.request)
		})
	}
}

func TestGopherMartHandler_Withdraw(t *testing.T) {
	type fields struct {
		db     DataBaser
		logger *zap.Logger
	}
	type args struct {
		writer  http.ResponseWriter
		request *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &GopherMartHandler{
				db:     tt.fields.db,
				logger: tt.fields.logger,
			}
			h.Withdraw(tt.args.writer, tt.args.request)
		})
	}
}

func TestGopherMartHandler_Withdrawals(t *testing.T) {
	type fields struct {
		db     DataBaser
		logger *zap.Logger
	}
	type args struct {
		writer  http.ResponseWriter
		request *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &GopherMartHandler{
				db:     tt.fields.db,
				logger: tt.fields.logger,
			}
			h.Withdrawals(tt.args.writer, tt.args.request)
		})
	}
}

func TestNewGopherMartHandler(t *testing.T) {
	type args struct {
		db     DataBaser
		logger *zap.Logger
	}
	tests := []struct {
		name string
		args args
		want *GopherMartHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGopherMartHandler(tt.args.db, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGopherMartHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
