package httpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"

	"github.com/h2p2f/loyalty-gophermart/internal/gophermart/models"
	"github.com/h2p2f/loyalty-gophermart/internal/gophermart/utils/jwt"
	"github.com/h2p2f/loyalty-gophermart/internal/gophermart/utils/luhn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

// GopherMartHandler is a struct for http handlers
type GopherMartHandler struct {
	db     DataBaser
	logger *zap.Logger
}

// NewGopherMartHandler is a function that returns a new GopherMartHandler
func NewGopherMartHandler(db DataBaser, logger *zap.Logger) *GopherMartHandler {
	return &GopherMartHandler{db: db, logger: logger}
}

// Register is a function that registers a new user
func (h *GopherMartHandler) Register(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		h.logger.Sugar().Errorf("Method not allowed: %s", request.Method)
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var buf bytes.Buffer
	_, err := buf.ReadFrom(request.Body)
	if err != nil {
		h.logger.Sugar().Errorf("Error reading request body: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	user := models.User{}
	err = json.Unmarshal(buf.Bytes(), &user)
	if err != nil {
		h.logger.Sugar().Errorf("Error unmarshalling request body: %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		h.logger.Sugar().Errorf("Error hashing password: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	user.Password = string(cryptedPassword)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err = h.db.NewUser(ctx, user)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			h.logger.Sugar().Errorf("Error creating user: %v", err)
			writer.WriteHeader(http.StatusConflict)
			return
		}
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.logger.Sugar().Infof("User %s created", user.Login)
	token, err := jwt.GenerateToken(user.Login)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.logger.Sugar().Infof("User %s logged in", user.Login)
	writer.Header().Add("Authorization", "Bearer "+token)
	writer.WriteHeader(http.StatusOK)

}

// Login is a function that logs in a user
func (h *GopherMartHandler) Login(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		h.logger.Sugar().Errorf("Method not allowed: %s", request.Method)
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var buf bytes.Buffer
	_, err := buf.ReadFrom(request.Body)
	if err != nil {
		h.logger.Sugar().Errorf("Error reading request body: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	user := models.User{}
	err = json.Unmarshal(buf.Bytes(), &user)
	if err != nil {
		h.logger.Sugar().Errorf("Error unmarshalling request body: %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	pass, err := h.db.FindPassByLogin(ctx, user.Login)
	if err != nil {
		h.logger.Sugar().Errorf("Error finding user: %v", err)
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(pass), []byte(user.Password))
	if err != nil {
		h.logger.Sugar().Errorf("Error comparing passwords: %v", err)
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	token, err := jwt.GenerateToken(user.Login)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.logger.Sugar().Infof("User %s logged in", user.Login)
	writer.Header().Add("Authorization", "Bearer "+token)
	writer.WriteHeader(http.StatusOK)
}

// AddOrder is a function that adds a new order
func (h *GopherMartHandler) AddOrder(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		h.logger.Sugar().Errorf("Method not allowed: %s", request.Method)
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	login, ok := LoginFromContext(request.Context())
	if !ok {
		h.logger.Sugar().Errorf("Error getting login from context: %v", ok)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.logger.Sugar().Infof("User %s is adding order", login)

	var buf bytes.Buffer
	_, err := buf.ReadFrom(request.Body)
	if err != nil {
		h.logger.Sugar().Errorf("Error reading request body: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	order := buf.String()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	owner, found := h.db.CheckUniqueOrder(ctx, order)
	if found {
		if owner == login {
			h.logger.Sugar().Info("order already exists")
			writer.WriteHeader(http.StatusOK)
			return
		}
		h.logger.Sugar().Info("order already exists")
		writer.WriteHeader(http.StatusConflict)
		return
	}
	if !luhn.Validate(order) {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	timeCreated := time.Now()
	orderModel := models.Order{
		Number:      order,
		Status:      models.NEW,
		Accrual:     0,
		TimeCreated: timeCreated,
	}
	err = h.db.NewOrder(ctx, login, orderModel)
	if err != nil {
		h.logger.Sugar().Errorf("Error creating order: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.logger.Sugar().Infof("Order %s created", order)
	writer.WriteHeader(http.StatusAccepted)

}

// Orders is a function that returns all orders of a user
func (h *GopherMartHandler) Orders(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	login, ok := LoginFromContext(request.Context())
	if !ok {
		h.logger.Sugar().Errorf("Error getting login from context: %v", ok)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	orders, err := h.db.GetOrdersByUser(ctx, login)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.logger.Sugar().Infof("User %s is getting orders", login)
	h.logger.Sugar().Infof("Orders: %v", string(orders))
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(orders)
	if err != nil {
		h.logger.Sugar().Error(err)
	}

}

// Balance is a function that returns the balance of a user
func (h *GopherMartHandler) Balance(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	login, ok := LoginFromContext(request.Context())
	if !ok {
		h.logger.Sugar().Errorf("Error getting login from context: %v", ok)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	balance, err := h.db.GetBalance(ctx, login)
	if err != nil {
		h.logger.Sugar().Errorf("Error getting balance: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	withdraws := h.db.GetSumOfAllWithdraws(ctx, login)
	account := models.Account{Balance: balance, Withdraws: withdraws}
	resp, err := json.Marshal(account)
	if err != nil {
		h.logger.Sugar().Errorf("Error marshalling account: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	_, err = writer.Write(resp)
	if err != nil {
		h.logger.Sugar().Error(err)
	}
}

// Withdraw is a function that withdraws money from a user's account
func (h *GopherMartHandler) Withdraw(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	login, ok := LoginFromContext(request.Context())
	if !ok {
		h.logger.Sugar().Errorf("Error getting login from context: %v", ok)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	_, err := buf.ReadFrom(request.Body)
	if err != nil {
		h.logger.Sugar().Errorf("Error reading request body: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	withdraw := models.Withdraw{}
	err = json.Unmarshal(buf.Bytes(), &withdraw)
	if err != nil {
		h.logger.Sugar().Errorf("Error unmarshalling withdraw: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if withdraw.Sum <= 0 {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	if !luhn.Validate(withdraw.Order) {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	balance, err := h.db.GetBalance(ctx, login)
	if err != nil {
		h.logger.Sugar().Errorf("Error getting balance: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if balance < withdraw.Sum {
		writer.WriteHeader(http.StatusPaymentRequired)
		return
	}
	orderModel := models.Order{
		Number:      withdraw.Order,
		Status:      models.PROCESSED,
		Accrual:     0,
		TimeCreated: time.Now(),
	}
	err = h.db.NewOrder(ctx, login, orderModel)
	if err != nil {
		h.logger.Sugar().Errorf("Error create new order %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	withdraw.TimeCreated = time.Now()
	err = h.db.NewWithdraw(ctx, login, withdraw)
	if err != nil {
		h.logger.Sugar().Errorf("Error create new withdraw %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)

}

// Withdrawals is a function that returns all withdraws of a user
func (h *GopherMartHandler) Withdrawals(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	login, ok := LoginFromContext(request.Context())
	if !ok {
		h.logger.Sugar().Errorf("Error getting login from context: %v", ok)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	withdraws := h.db.GetAllWithdraws(ctx, login)
	if withdraws == nil {
		writer.WriteHeader(http.StatusNoContent)
		return
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	_, err := writer.Write(withdraws)
	if err != nil {
		h.logger.Sugar().Error(err)
	}
}
