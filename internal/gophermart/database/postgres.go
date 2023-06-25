package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/h2p2f/loyalty-gophermart/internal/gophermart/models"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

type PostgresDB struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewPostgresDB creates new PostgresDB
func NewPostgresDB(param string, logger *zap.Logger) *PostgresDB {
	db, err := sql.Open("pgx", param)
	if err != nil {
		logger.Sugar().Errorf("Error opening database connection: %v", err)
	}
	return &PostgresDB{db: db, logger: logger}
}

// Close database connection
func (pgdb *PostgresDB) Close() {
	err := pgdb.db.Close()
	if err != nil {
		pgdb.logger.Sugar().Errorf("Error closing database connection: %v", err)
	}
}

// Create creates tables in database
func (pgdb *PostgresDB) Create() error {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	query := `CREATE TABLE IF NOT EXISTS go_mart_user (
			login text not null PRIMARY KEY,
			password text not null
	);`

	_, err := pgdb.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	query = `CREATE TABLE IF NOT EXISTS go_mart_order (
    		id text not null PRIMARY KEY,
    		uuid text not null REFERENCES go_mart_user(login),
    		status text,
    		accrual numeric(10,2) not null,
    		time_created timestamp not null
    		);`
	_, err = pgdb.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	query = `CREATE TABLE IF NOT EXISTS go_mart_user_balance (
    		uuid text not null PRIMARY KEY REFERENCES go_mart_user(login),
    		balance numeric(10,2) not null
	);`
	_, err = pgdb.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	query = `CREATE TABLE IF NOT EXISTS go_mart_withdraws (
    		uuid text not null REFERENCES go_mart_user(login),
    		order_id text not null REFERENCES go_mart_order(id),
    		amount numeric(10,2) not null,
    		time_created timestamp not null
    		    	);`
	_, err = pgdb.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

// NewUser creates new user in database
func (pgdb *PostgresDB) NewUser(ctx context.Context, user models.User) error {
	query := `INSERT INTO go_mart_user
   			(login, password)
				VALUES ($1, $2)`
	_, err := pgdb.db.ExecContext(ctx, query, user.Login, user.Password)
	if err != nil {
		return err
	}
	query = `INSERT INTO go_mart_user_balance
       			(uuid, balance)
       				VALUES ($1, $2)`
	_, err = pgdb.db.ExecContext(ctx, query, user.Login, 0)
	if err != nil {
		return err
	}
	return nil
}

// FindPassByLogin finds password by login in database
func (pgdb *PostgresDB) FindPassByLogin(ctx context.Context, login string) (string, error) {
	var password string
	query := `SELECT password FROM go_mart_user WHERE login = $1`
	err := pgdb.db.QueryRowContext(ctx, query, login).Scan(&password)
	if err != nil {
		return "", err
	}
	return password, nil
}

// NewOrder creates new order in database
func (pgdb *PostgresDB) NewOrder(ctx context.Context, login string, order models.Order) error {
	query := `INSERT INTO go_mart_order
       			(id, uuid, status, accrual, time_created)
       				VALUES ($1, $2, $3, $4, $5)`
	_, err := pgdb.db.ExecContext(ctx, query, order.Number, login, order.Status, order.Accrual, order.TimeCreated)
	if err != nil {
		return err
	}
	return nil
}

// GetOrdersByUser gets orders by user from database
func (pgdb *PostgresDB) GetOrdersByUser(ctx context.Context, login string) ([]byte, error) {
	var ordersResult []byte
	query := `SELECT id, status, accrual, time_created FROM go_mart_order WHERE uuid = $1 ORDER BY time_created`
	rows, err := pgdb.db.QueryContext(ctx, query, login)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, err
	}
	defer rows.Close()
	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err = rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.TimeCreated)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	if len(orders) == 0 {
		return nil, errors.New("no orders found")
	}
	ordersResult, err = json.Marshal(orders)
	if err != nil {
		return nil, err
	}
	return ordersResult, nil

}

// CheckUniqueOrder checks if order is unique
func (pgdb *PostgresDB) CheckUniqueOrder(ctx context.Context, order string) (string, error) {
	var st string
	query := `SELECT uuid FROM go_mart_order WHERE id = $1`
	err := pgdb.db.QueryRowContext(ctx, query, order).Scan(&st)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	return st, nil
}

// GetBalance gets balance from database
func (pgdb *PostgresDB) GetBalance(ctx context.Context, login string) (float64, error) {
	var balance float64
	query := `SELECT balance FROM go_mart_user_balance WHERE uuid = $1`
	err := pgdb.db.QueryRowContext(ctx, query, login).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

// GetSumOfAllWithdraws gets sum of all withdraws from database
func (pgdb *PostgresDB) GetSumOfAllWithdraws(ctx context.Context, login string) float64 {
	var withdraws float64
	query := `SELECT SUM(amount) FROM go_mart_withdraws WHERE uuid = $1`
	err := pgdb.db.QueryRowContext(ctx, query, login).Scan(&withdraws)
	if err != nil {
		return 0
	}

	return withdraws
}

// GetAllWithdraws gets all withdraws from database
func (pgdb *PostgresDB) GetAllWithdraws(ctx context.Context, login string) []byte {
	var withdrawsResult []byte
	query := `SELECT order_id, amount, time_created FROM go_mart_withdraws WHERE uuid = $1`

	rows, err := pgdb.db.QueryContext(ctx, query, login)
	if err != nil {
		return nil
	}
	if rows.Err() != nil {
		return nil
	}
	defer rows.Close()
	var withdraws []models.Withdraw
	for rows.Next() {
		var withdraw models.Withdraw
		err = rows.Scan(&withdraw.Order, &withdraw.Sum, &withdraw.TimeCreated)
		if err != nil {
			return nil
		}
		withdraws = append(withdraws, withdraw)
	}
	withdrawsResult, err = json.Marshal(withdraws)
	if err != nil {
		return nil
	}
	return withdrawsResult
}

// NewWithdraw creates new withdraw in database
func (pgdb *PostgresDB) NewWithdraw(ctx context.Context, login string, withdraw models.Withdraw) error {
	query := `INSERT INTO go_mart_withdraws
       			(uuid, order_id, amount, time_created)
       				VALUES ($1, $2, $3, $4)`
	_, err := pgdb.db.ExecContext(ctx, query, login, withdraw.Order, withdraw.Sum, withdraw.TimeCreated)
	if err != nil {
		return err
	}
	query = `UPDATE go_mart_user_balance SET balance = balance - $1 WHERE uuid = $2`
	_, err = pgdb.db.ExecContext(ctx, query, withdraw.Sum, login)
	if err != nil {
		return err
	}
	return nil
}

// GetUnfinishedOrders gets unfinished orders from database
func (pgdb *PostgresDB) GetUnfinishedOrders() (map[string]string, error) {
	orders := make(map[string]string)
	rows, err := pgdb.db.Query(
		`SELECT id, status FROM go_mart_order WHERE status = $1 OR status = $2`, models.NEW, models.PROCESSING)
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var order, status string
		err = rows.Scan(&order, &status)
		if err != nil {
			return nil, err
		}
		orders[order] = status
	}
	return orders, nil
}

// UpdateOrderStatus updates order status in database
func (pgdb *PostgresDB) UpdateOrderStatus(order, status string, accrual float64) error {
	tx, err := pgdb.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE go_mart_order SET status = $1, accrual = $2 WHERE id = $3`
	_, err = pgdb.db.Exec(query, status, accrual, order)
	if err != nil {
		pgdb.logger.Sugar().Errorf("Error, need rollback transaction: %v", err)
		return err
	}
	query = `UPDATE go_mart_user_balance SET balance = balance + $1 WHERE uuid = (SELECT uuid FROM go_mart_order WHERE id = $2)`
	_, err = pgdb.db.Exec(query, accrual, order)
	if err != nil {
		pgdb.logger.Sugar().Errorf("Error, need rollback transaction: %v", err)
		return err
	}
	err = tx.Commit()
	if err != nil {
		pgdb.logger.Sugar().Errorf("Error commit transaction: %v", err)
	}
	return nil
}
