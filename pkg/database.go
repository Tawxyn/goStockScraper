package database

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgres struct {
	db *pgxpool.Pool
}

// StockCashFlow struct
type StockCashFlow struct {
	Ticker       string
	CashFlow2020 int
	CashFlow2021 int
	CashFlow2022 int
	CashFlow2023 int
}

var (
	pgInstance *postgres
	pgOnce     sync.Once
)

// Initalize the database with pgxpool
func InitDatabase(ctx context.Context, connString string) (*postgres, error) {
	var err error

	pgOnce.Do(func() {
		var db *pgxpool.Pool
		db, err = pgxpool.New(ctx, connString)
		if err == nil {
			pgInstance = &postgres{db}
		}
	})

	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	return pgInstance, nil
}

func (pg *postgres) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

// Close function to shutdown gracefully
func (pg *postgres) Close() {
	pg.db.Close()
}

// InsertFCF inserts cash flow values for 2020-2023 years
func (pg *postgres) InsertFCF(ctx context.Context, ticker string, cashFlow2020, cashFlow2021, cashFlow2022, cashFlow2023 string) error {
	var err error

	// Convert string values to integers
	cf2020, err := strconv.Atoi(cashFlow2020)
	if err != nil {
		return fmt.Errorf("failed to convert cashFlow2020 to integer: %v", err)
	}

	cf2021, err := strconv.Atoi(cashFlow2021)
	if err != nil {
		return fmt.Errorf("failed to convert cashFlow2021 to integer: %v", err)
	}

	cf2022, err := strconv.Atoi(cashFlow2022)
	if err != nil {
		return fmt.Errorf("failed to convert cashFlow2022 to integer: %v", err)
	}

	cf2023, err := strconv.Atoi(cashFlow2023)
	if err != nil {
		return fmt.Errorf("failed to convert cashFlow2023 to integer: %v", err)
	}
	query := `
        INSERT INTO stock_cash_flow 
            (ticker, cash_flow_2020, cash_flow_2021, cash_flow_2022, cash_flow_2023) 
        VALUES 
            ($1, $2, $3, $4, $5)`

	_, err = pg.db.Exec(ctx, query, ticker, cf2020, cf2021, cf2022, cf2023)
	if err != nil {
		return fmt.Errorf("failed to insert data into database: %v", err)
	}

	return nil
}

func (pg *postgres) CheckTickerExists(ctx context.Context, ticker string) (bool, error) {
	var count int
	query := `SELECT COUNT(*)
			  FROM stock_cash_flow
			  WHERE ticker =$1`

	// QueryRowContext executes a query that is expected to return at most one row.
	err := pg.db.QueryRow(ctx, query, ticker).Scan(&count)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil // No row found for the ticker
		}
		return false, fmt.Errorf("error checking ticker: %v", err)
	}

	// If count > 0, the ticker exists; otherwise, it does not
	return count > 0, nil
}
