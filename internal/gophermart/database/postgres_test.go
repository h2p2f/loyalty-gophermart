package database

import (
	"context"
	"database/sql"
	"errors"
	"github.com/h2p2f/loyalty-gophermart/internal/gophermart/models"
	"go.uber.org/zap"
	"log"
	"reflect"
	"testing"
	"time"
)

func TestPostgresDB_CheckUniqueOrder(t *testing.T) {
	db, err := sql.Open("pgx", "postgres://practicum:yandex@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	logger := zap.NewExample()
	pg := &PostgresDB{
		db:     db,
		logger: logger,
	}
	ctx := context.Background()

	tests := []struct {
		name  string
		order string
		user  string
		want  error
	}{
		{
			name:  "positive test",
			order: "12345678903",
			user:  "FirstUser",
			want:  nil,
		},
		{
			name:  "negative test",
			order: "1234567890",
			user:  "",
			want:  errors.New("sql: no rows in result set"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := pg.CheckUniqueOrder(ctx, tt.order)

			if got != tt.user {
				t.Errorf("CheckUniqueOrder() got = %v, want %v", got, tt.user)
			}
			if tt.want != nil && errors.Is(got1, tt.want) {
				t.Errorf("CheckUniqueOrder() got1 = %v, want %v", got1, tt.want)
			}
			//if !reflect.DeepEqual(got1, tt.want) {
			//	t.Errorf("CheckUniqueOrder() got1 = %v, want %v", got1, tt.want)
			//}

		})
	}
}

func TestPostgresDB_FindPassByLogin(t *testing.T) {
	db, err := sql.Open("pgx", "postgres://practicum:yandex@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	logger := zap.NewExample()
	pg := &PostgresDB{
		db:     db,
		logger: logger,
	}
	ctx := context.Background()

	tests := []struct {
		name    string
		user    string
		want    string
		wantErr bool
	}{
		{
			name:    "positive test",
			user:    "FirstUser",
			want:    "$2a$10$4s6.ghWw25/q2fxwLNh/N.UVMDNTK/GhQNR9P2JZALP.bX97ttwOe",
			wantErr: false,
		},
		{
			name:    "negative test",
			user:    "SecondUser",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := pg.FindPassByLogin(ctx, tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindPassByLogin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FindPassByLogin() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostgresDB_GetAllWithdraws(t *testing.T) {
	db, err := sql.Open("pgx", "postgres://practicum:yandex@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	logger := zap.NewExample()
	pg := &PostgresDB{
		db:     db,
		logger: logger,
	}
	ctx := context.Background()

	tests := []struct {
		name string
		user string
		want []byte
	}{
		{
			name: "positive test",
			user: "FirstUser",
			want: []byte(`null`),
		},
		{
			name: "negative test",
			user: "SecondUser",
			want: []byte(`null`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := pg.GetAllWithdraws(ctx, tt.user); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllWithdraws() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostgresDB_GetBalance(t *testing.T) {
	db, err := sql.Open("pgx", "postgres://practicum:yandex@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	logger := zap.NewExample()
	pg := &PostgresDB{
		db:     db,
		logger: logger,
	}
	ctx := context.Background()
	tests := []struct {
		name    string
		user    string
		want    float64
		wantErr bool
	}{
		{
			name:    "positive test",
			user:    "FirstUser",
			want:    0,
			wantErr: false,
		},
		{
			name:    "negative test",
			user:    "SecondUser",
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := pg.GetBalance(ctx, tt.user)
			if (err != nil) != tt.wantErr {

				t.Errorf("GetBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetBalance() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostgresDB_GetOrdersByUser(t *testing.T) {
	db, err := sql.Open("pgx", "postgres://practicum:yandex@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	logger := zap.NewExample()
	pg := &PostgresDB{
		db:     db,
		logger: logger,
	}
	ctx := context.Background()

	tests := []struct {
		name    string
		user    string
		want    []byte
		wantErr bool
	}{
		{
			name:    "positive test",
			user:    "FirstUser",
			want:    []byte(`[{"number":"12345678903","status":"NEW","accrual":0,"uploaded_at":"2023-06-09T14:31:04.222794Z"},{"number":"7374867609","status":"NEW","accrual":0,"uploaded_at":"2023-06-09T14:31:24.155728Z"}]`),
			wantErr: false,
		},
		{
			name:    "negative test",
			user:    "SecondUser",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := pg.GetOrdersByUser(ctx, tt.user)
			//this is for debug
			//fmt.Println(string(got))
			if (err != nil) != tt.wantErr {

				t.Errorf("GetOrdersByUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetOrdersByUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostgresDB_GetSumOfAllWithdraws(t *testing.T) {
	db, err := sql.Open("pgx", "postgres://practicum:yandex@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	logger := zap.NewExample()
	pg := &PostgresDB{
		db:     db,
		logger: logger,
	}
	ctx := context.Background()

	tests := []struct {
		name string
		user string
		want float64
	}{
		{
			name: "positive test",
			user: "FirstUser",
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := pg.GetSumOfAllWithdraws(ctx, tt.user); got != tt.want {
				t.Errorf("GetSumOfAllWithdraws() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostgresDB_GetUnfinishedOrders(t *testing.T) {
	db, err := sql.Open("pgx", "postgres://practicum:yandex@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	logger := zap.NewExample()
	pg := &PostgresDB{
		db:     db,
		logger: logger,
	}
	tests := []struct {
		name    string
		want    map[string]string
		wantErr bool
	}{
		{
			name: "positive test",
			want: map[string]string{
				"12345678903": "NEW",
				"7374867609":  "NEW",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := pg.GetUnfinishedOrders()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUnfinishedOrders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUnfinishedOrders() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostgresDB_NewOrder(t *testing.T) {
	db, err := sql.Open("pgx", "postgres://practicum:yandex@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	logger := zap.NewExample()
	pg := &PostgresDB{
		db:     db,
		logger: logger,
	}
	ctx := context.Background()

	tests := []struct {
		name             string
		user             string
		orderNumber      string
		orderStatus      string
		orderAccrual     float64
		orderTimeCreated time.Time
		order            models.Order
		wantErr          bool
	}{
		{
			name:             "positive test",
			user:             "FirstUser",
			orderNumber:      "1234567891",
			orderStatus:      "NEW",
			orderAccrual:     10.01,
			orderTimeCreated: time.Date(2023, 6, 9, 14, 31, 4, 222794000, time.UTC),
			wantErr:          false,
		},
		{
			name:             "negative test",
			user:             "SecondUser",
			orderNumber:      "1234567892",
			orderStatus:      "NEW",
			orderAccrual:     10.01,
			orderTimeCreated: time.Date(2023, 6, 9, 14, 31, 4, 222794000, time.UTC),
			wantErr:          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := models.Order{
				Number:      tt.orderNumber,
				Status:      tt.orderStatus,
				Accrual:     &tt.orderAccrual,
				TimeCreated: tt.orderTimeCreated,
			}
			if err := pg.NewOrder(ctx, tt.user, order); (err != nil) != tt.wantErr {
				t.Errorf("NewOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr == false {
				_, err := pg.db.Exec("DELETE FROM go_mart_order WHERE id = $1", tt.order.Number)
				if err != nil {
					t.Errorf("can't delete order: %v", err)
				}
			}
		})
	}
}

func TestPostgresDB_NewUser(t *testing.T) {
	db, err := sql.Open("pgx", "postgres://practicum:yandex@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	logger := zap.NewExample()
	pg := &PostgresDB{
		db:     db,
		logger: logger,
	}
	ctx := context.Background()

	tests := []struct {
		name    string
		user    models.User
		wantErr bool
	}{
		{
			name: "positive test",
			user: models.User{
				Login:    "TestUser",
				Password: "TestPassword",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := pg.NewUser(ctx, tt.user); (err != nil) != tt.wantErr {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr == false {
				_, err := pg.db.Exec("DELETE FROM go_mart_user_balance WHERE uuid = $1", tt.user.Login)
				if err != nil {
					t.Errorf("can't delete user: %v", err)
				}
				_, _ = pg.db.Exec("DELETE FROM go_mart_user WHERE login = $1", tt.user.Login)

			}
		})
	}
}

func TestPostgresDB_NewWithdraw(t *testing.T) {
	db, err := sql.Open("pgx", "postgres://practicum:yandex@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	logger := zap.NewExample()
	pg := &PostgresDB{
		db:     db,
		logger: logger,
	}
	ctx := context.Background()

	tests := []struct {
		name         string
		user         string
		withdraw     models.Withdraw
		wSum         float64
		order        models.Order
		oNumber      string
		oStatus      string
		oAccrual     float64
		oTimeCreated time.Time
		wantErr      bool
	}{
		{
			name:         "positive test",
			user:         "FirstUser",
			wSum:         100,
			oNumber:      "12345678912",
			oStatus:      "NEW",
			oAccrual:     0,
			oTimeCreated: time.Date(2023, 6, 9, 14, 31, 4, 222794000, time.UTC),
			wantErr:      false,
		},
		{
			name:         "negative test",
			user:         "SecondUser",
			wSum:         100,
			oNumber:      "1234567892",
			oStatus:      "NEW",
			oAccrual:     0,
			oTimeCreated: time.Date(2023, 6, 9, 14, 31, 4, 222794000, time.UTC),
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withdraw := models.Withdraw{
				Order:       tt.oNumber,
				Sum:         tt.wSum,
				TimeCreated: tt.oTimeCreated,
			}
			order := models.Order{
				Number:      tt.oNumber,
				Status:      tt.oStatus,
				Accrual:     &tt.oAccrual,
				TimeCreated: tt.oTimeCreated,
			}
			if !tt.wantErr {
				if err := pg.NewOrder(ctx, tt.user, order); (err != nil) != tt.wantErr {
					t.Errorf("NewOrder() error = %v, wantErr %v", err, tt.wantErr)
				}
				_, err := pg.db.Exec("UPDATE go_mart_user_balance SET balance = 1000 WHERE uuid = $1", tt.user)
				if err != nil {
					t.Errorf("can't update balance: %v", err)
				}
			}
			if err := pg.NewWithdraw(ctx, tt.user, withdraw); (err != nil) != tt.wantErr {
				t.Errorf("NewWithdraw() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				_, err = pg.db.Exec("DELETE FROM go_mart_withdraws WHERE order_id = $1", withdraw.Order)
				if err != nil {
					t.Errorf("can't delete withdraw: %v", err)
				}
				_, err := pg.db.Exec("DELETE FROM go_mart_order WHERE id = $1", order.Number)
				if err != nil {
					t.Errorf("can't delete order: %v", err)
				}

				_, err = pg.db.Exec("UPDATE go_mart_user_balance SET balance = 0 WHERE uuid = $1", tt.user)
				if err != nil {
					t.Errorf("can't update balance: %v", err)
				}
			}

		})
	}
}

func TestPostgresDB_UpdateOrderStatus(t *testing.T) {
	db, err := sql.Open("pgx", "postgres://practicum:yandex@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()
	logger := zap.NewExample()
	pg := &PostgresDB{
		db:     db,
		logger: logger,
	}

	tests := []struct {
		name    string
		order   string
		status  string
		accrual float64
		wantErr bool
	}{
		{
			name:    "positive test",
			order:   "12345678903",
			status:  "NEW",
			accrual: 0,
			wantErr: false,
		},
		{
			name:    "negative test",
			order:   "12345678904",
			status:  "NEW",
			accrual: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := pg.UpdateOrderStatus(tt.order, tt.status, tt.accrual); (err != nil) != tt.wantErr {
				t.Errorf("UpdateOrderStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
