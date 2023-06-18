package database

import (
	"context"
	"database/sql"
	"github.com/h2p2f/loyalty-gophermart/internal/gophermart/models"
	"go.uber.org/zap"
	"log"
	"reflect"
	"testing"
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
		want  bool
	}{
		{
			name:  "positive test",
			order: "12345678903",
			user:  "FirstUser",
			want:  true,
		},
		{
			name:  "negative test",
			order: "1234567890",
			user:  "",
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := pg.CheckUniqueOrder(ctx, tt.order)
			if got != tt.user {
				t.Errorf("CheckUniqueOrder() got = %v, want %v", got, tt.user)
			}
			if got1 != tt.want {
				t.Errorf("CheckUniqueOrder() got1 = %v, want %v", got1, tt.want)
			}
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
			name: "positive test",
			user: "FirstUser",

			want:    []byte(`[{"number":"12345678903","status":"NEW","uploaded_at":"2023-06-09T14:31:04.222794Z"},{"number":"7374867609","status":"NEW","uploaded_at":"2023-06-09T14:31:24.155728Z"}]`),
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

//func TestPostgresDB_NewOrder(t *testing.T) {
//	db, err := sql.Open("pgx", "postgres://practicum:yandex@localhost:5432/postgres?sslmode=disable")
//	if err != nil {
//		log.Fatal(err)
//		return
//	}
//	defer func() {
//		if err := db.Close(); err != nil {
//			t.Errorf("can't close db: %v", err)
//		}
//	}()
//	logger := zap.NewExample()
//	pg := &PostgresDB{
//		db:     db,
//		logger: logger,
//	}
//	ctx := context.Background()
//
//	tests := []struct {
//		name    string
//		user    string
//		order   models.Order
//		wantErr bool
//	}{
//		{
//			name: "positive test",
//			user: "FirstUser",
//			order: models.Order{
//				Number:      "12345678903",
//				Status:      "NEW",
//				TimeCreated: time.Date(2023, 6, 9, 14, 31, 4, 222794000, time.UTC),
//			},
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//
//			if err := pg.NewOrder(ctx, tt.user, tt.order); (err != nil) != tt.wantErr {
//				t.Errorf("NewOrder() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}

func TestPostgresDB_NewUser(t *testing.T) {
	type fields struct {
		db     *sql.DB
		logger *zap.Logger
	}
	type args struct {
		ctx  context.Context
		user models.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pgdb := &PostgresDB{
				db:     tt.fields.db,
				logger: tt.fields.logger,
			}
			if err := pgdb.NewUser(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPostgresDB_NewWithdraw(t *testing.T) {
	type fields struct {
		db     *sql.DB
		logger *zap.Logger
	}
	type args struct {
		ctx      context.Context
		login    string
		withdraw models.Withdraw
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pgdb := &PostgresDB{
				db:     tt.fields.db,
				logger: tt.fields.logger,
			}
			if err := pgdb.NewWithdraw(tt.args.ctx, tt.args.login, tt.args.withdraw); (err != nil) != tt.wantErr {
				t.Errorf("NewWithdraw() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPostgresDB_UpdateOrderStatus(t *testing.T) {
	type fields struct {
		db     *sql.DB
		logger *zap.Logger
	}
	type args struct {
		order   string
		status  string
		accrual float64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pgdb := &PostgresDB{
				db:     tt.fields.db,
				logger: tt.fields.logger,
			}
			if err := pgdb.UpdateOrderStatus(tt.args.order, tt.args.status, tt.args.accrual); (err != nil) != tt.wantErr {
				t.Errorf("UpdateOrderStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
