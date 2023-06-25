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
	// Check if method is POST
	if request.Method != http.MethodPost {
		h.logger.Sugar().Errorf("Method not allowed: %s", request.Method)
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	key, ok := KeyFromContext(request.Context())
	if !ok {
		h.logger.Sugar().Errorf("Error getting key from context")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Read request body
	var buf bytes.Buffer
	_, err := buf.ReadFrom(request.Body)
	if err != nil {
		h.logger.Sugar().Errorf("Error reading request body: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Unmarshal request body
	user := models.User{}
	err = json.Unmarshal(buf.Bytes(), &user)
	if err != nil {
		h.logger.Sugar().Errorf("Error unmarshalling request body: %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	// Prepare password to be hashed
	cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		h.logger.Sugar().Errorf("Error hashing password: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	user.Password = string(cryptedPassword)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	//create user in database
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
	// Generate JWT token
	token, err := jwt.GenerateToken(user.Login, key)
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
	// Check if method is POST
	if request.Method != http.MethodPost {
		h.logger.Sugar().Errorf("Method not allowed: %s", request.Method)
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	key, ok := KeyFromContext(request.Context())
	if !ok {
		h.logger.Sugar().Errorf("Error getting key from context")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Read request body
	var buf bytes.Buffer
	_, err := buf.ReadFrom(request.Body)
	if err != nil {
		h.logger.Sugar().Errorf("Error reading request body: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Unmarshal request body
	user := models.User{}
	err = json.Unmarshal(buf.Bytes(), &user)
	if err != nil {
		h.logger.Sugar().Errorf("Error unmarshalling request body: %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	// Get password from database
	pass, err := h.db.FindPassByLogin(ctx, user.Login)
	if err != nil {
		h.logger.Sugar().Errorf("Error finding user: %v", err)
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(pass), []byte(user.Password))
	if err != nil {
		h.logger.Sugar().Errorf("Error comparing passwords: %v", err)
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	// Generate JWT token
	token, err := jwt.GenerateToken(user.Login, key)
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
	// Check if method is POST
	if request.Method != http.MethodPost {
		h.logger.Sugar().Errorf("Method not allowed: %s", request.Method)
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// Get login from context
	login, ok := LoginFromContext(request.Context())
	if !ok {
		h.logger.Sugar().Errorf("Error getting login from context: %v", ok)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.logger.Sugar().Infof("User %s is adding order", login)
	// Read request body
	var buf bytes.Buffer
	_, err := buf.ReadFrom(request.Body)
	if err != nil {
		h.logger.Sugar().Errorf("Error reading request body: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	// get order from request body
	order := buf.String()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	// check if order already exists
	owner, err := h.db.CheckUniqueOrder(ctx, order)
	if err != nil {
		h.logger.Sugar().Errorf("Error checking if order exists: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if owner != "" {
		if owner == login {
			h.logger.Sugar().Info("order already exists")
			writer.WriteHeader(http.StatusOK)
			return
		}
		//your password is incorrect, but the same password is used by the user:Jack17005 :)
		h.logger.Sugar().Info("this order's number conflict with another user's order")
		writer.WriteHeader(http.StatusConflict)
		return
	}
	// check if order's number is valid
	if !luhn.Validate(order) {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	timeCreated := time.Now()
	// create order in database
	var acc float64 = 0
	orderModel := models.Order{
		Number:      order,
		Status:      models.NEW,
		Accrual:     &acc,
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

// GetOrders is a function that returns all orders of a user
func (h *GopherMartHandler) GetOrders(writer http.ResponseWriter, request *http.Request) {
	// Check if method is GET
	if request.Method != http.MethodGet {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// Get login from context
	login, ok := LoginFromContext(request.Context())
	if !ok {
		h.logger.Sugar().Errorf("Error getting login from context: %v", ok)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	// Get orders from database
	orders, err := h.db.GetOrdersByUser(ctx, login)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.logger.Sugar().Infof("User %s is getting orders", login)
	h.logger.Sugar().Infof("GetOrders: %v", string(orders))
	// Write orders to response
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(orders)
	if err != nil {
		h.logger.Sugar().Error(err)
	}

}

// GetBalance is a function that returns the balance of a user
func (h *GopherMartHandler) GetBalance(writer http.ResponseWriter, request *http.Request) {
	// Check if method is GET
	if request.Method != http.MethodGet {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// Get login from context
	login, ok := LoginFromContext(request.Context())
	if !ok {
		h.logger.Sugar().Errorf("Error getting login from context: %v", ok)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	// Get balance from database
	balance, err := h.db.GetBalance(ctx, login)
	if err != nil {
		h.logger.Sugar().Errorf("Error getting balance: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Get sum of all withdraws from database
	withdraws := h.db.GetSumOfAllWithdraws(ctx, login)
	account := models.Account{Balance: balance, Withdraws: &withdraws}
	// Marshal account to json
	resp, err := json.Marshal(account)
	if err != nil {
		h.logger.Sugar().Errorf("Error marshalling account: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Write account to response
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	_, err = writer.Write(resp)
	if err != nil {
		h.logger.Sugar().Error(err)
	}
}

// DoWithdraw is a function that withdraws money from a user's account
func (h *GopherMartHandler) DoWithdraw(writer http.ResponseWriter, request *http.Request) {
	// Check if method is POST
	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// Get login from context
	login, ok := LoginFromContext(request.Context())
	if !ok {
		h.logger.Sugar().Errorf("Error getting login from context: %v", ok)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Read request body
	var buf bytes.Buffer
	_, err := buf.ReadFrom(request.Body)
	if err != nil {
		h.logger.Sugar().Errorf("Error reading request body: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Unmarshal request body to withdraw
	withdraw := models.Withdraw{}
	err = json.Unmarshal(buf.Bytes(), &withdraw)
	if err != nil {
		h.logger.Sugar().Errorf("Error unmarshalling withdraw: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Check if sum is valid
	if withdraw.Sum <= 0 {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	// Check if order is valid
	if !luhn.Validate(withdraw.Order) {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	// Get balance from database
	balance, err := h.db.GetBalance(ctx, login)
	if err != nil {
		h.logger.Sugar().Errorf("Error getting balance: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Check if balance is enough
	if balance < withdraw.Sum {
		writer.WriteHeader(http.StatusPaymentRequired)
		return
	}
	// Update balance and create order in database
	var acc float64 = 0
	orderModel := models.Order{
		Number:      withdraw.Order,
		Status:      models.PROCESSED,
		Accrual:     &acc,
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

// GetWithdrawals is a function that returns all withdraws of a user
func (h *GopherMartHandler) GetWithdrawals(writer http.ResponseWriter, request *http.Request) {
	// Check if method is GET
	if request.Method != http.MethodGet {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// Get login from context
	login, ok := LoginFromContext(request.Context())
	if !ok {
		h.logger.Sugar().Errorf("Error getting login from context: %v", ok)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	// Get all withdraws from database
	withdraws := h.db.GetAllWithdraws(ctx, login)
	if withdraws == nil {
		writer.WriteHeader(http.StatusNoContent)
		return
	}
	// Marshal withdraws to json
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	_, err := writer.Write(withdraws)
	if err != nil {
		h.logger.Sugar().Error(err)
	}
}
