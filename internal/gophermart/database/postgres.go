package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/h2p2f/loyalty-gophermart/internal/gophermart/models"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"time"
)

type PostgresDB struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewPostgresDB(param string, logger *zap.Logger) *PostgresDB {
	db, err := sql.Open("pgx", param)
	if err != nil {
		logger.Sugar().Errorf("Error opening database connection: %v", err)
	}
	return &PostgresDB{db: db, logger: logger}
}

func (pgdb *PostgresDB) Close() {
	err := pgdb.db.Close()
	if err != nil {
		pgdb.logger.Sugar().Errorf("Error closing database connection: %v", err)
	}
}

func (pgdb *PostgresDB) Create() error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

func (pgdb *PostgresDB) NewUser(login, password string) error {
	_, err := pgdb.db.Exec(
		`INSERT INTO go_mart_user
   			(login, password)
				VALUES ($1, $2)`,
		login, password)
	if err != nil {
		return err
	}
	_, err = pgdb.db.Exec(
		`INSERT INTO go_mart_user_balance
       			(uuid, balance)
       				VALUES ($1, $2)`, login, 0)
	if err != nil {
		return err
	}
	return nil
}

func (pgdb *PostgresDB) FindPassByLogin(login string) (string, error) {
	var password string
	err := pgdb.db.QueryRow(
		`SELECT password FROM go_mart_user WHERE login = $1`, login).Scan(&password)
	if err != nil {
		return "", err
	}
	return password, nil
}

func (pgdb *PostgresDB) NewOrder(id, login, status string, accrual float64, timeCreated time.Time) error {
	_, err := pgdb.db.Exec(
		`INSERT INTO go_mart_order
   			(id, uuid, status, accrual, time_created)
				VALUES ($1, $2, $3, $4, $5)`,
		id, login, status, accrual, timeCreated)
	if err != nil {
		return err
	}
	return nil
}

func (pgdb *PostgresDB) GetOrdersByUser(login string) ([]byte, error) {
	var ordersResult []byte
	rows, err := pgdb.db.Query(
		`SELECT 
    				id, status, accrual, time_created 
				FROM
				    go_mart_order 
				WHERE 
				    uuid = $1 
				ORDER BY 
				    time_created`,
		login)
	if err != nil {
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
	ordersResult, err = json.Marshal(orders)
	fmt.Println(string(ordersResult))
	if err != nil {
		return nil, err
	}
	return ordersResult, nil

}

func (pgdb *PostgresDB) CheckUniqueOrder1(order string) bool {
	var st bool
	err := pgdb.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM go_mart_order WHERE id = $1)`, order).Scan(&st)
	if err != nil {
		return false
	}
	return st
}

func (pgdb *PostgresDB) CheckUniqueOrder(order string) (string, bool) {
	var st string
	err := pgdb.db.QueryRow(
		`SELECT uuid FROM go_mart_order WHERE id = $1`, order).Scan(&st)
	if err != nil {
		return "", false
	}
	if st == "" {
		return "", false
	}
	return st, true
}

func (pgdb *PostgresDB) GetBalance(login string) (float64, error) {
	var balance float64
	err := pgdb.db.QueryRow(
		`SELECT balance FROM go_mart_user_balance WHERE uuid = $1`, login).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

func (pgdb *PostgresDB) GetSumOfAllWithdraws(login string) float64 {
	var withdraws float64
	err := pgdb.db.QueryRow(
		`SELECT SUM(amount) FROM go_mart_withdraws WHERE uuid = $1`, login).Scan(&withdraws)
	if err != nil {
		return 0
	}
	return withdraws
}

func (pgdb *PostgresDB) GetAllWithdraws(login string) []byte {
	var withdrawsResult []byte
	rows, err := pgdb.db.Query(
		`SELECT order_id, amount, time_created FROM go_mart_withdraws WHERE uuid = $1`, login)
	if err != nil {
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

func (pgdb *PostgresDB) NewWithdraw(login, order string, amount float64, timeCreated time.Time) error {
	_, err := pgdb.db.Exec(
		`INSERT INTO go_mart_withdraws
   			(uuid, order_id, amount, time_created)
				VALUES ($1, $2, $3, $4)`,
		login, order, amount, timeCreated)
	if err != nil {
		return err
	}
	return nil
}

func (pgdb *PostgresDB) GetUnfinishedOrders() (map[string]string, error) {
	orders := make(map[string]string)
	rows, err := pgdb.db.Query(
		`SELECT id, status FROM go_mart_order WHERE status = $1 OR status = $2`, models.NEW, models.PROCESSING)
	if err != nil {
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

func (pgdb *PostgresDB) UpdateOrderStatus(order, status string, accrual float64) error {
	_, err := pgdb.db.Exec(
		`UPDATE go_mart_order SET status = $1, accrual = $2 WHERE id = $3`,
		status, accrual, order)
	if err != nil {
		return err
	}
	return nil
}
