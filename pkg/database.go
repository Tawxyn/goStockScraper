package database

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Create a single instance
var pool *pgxpool.Pool

// Initalize the database with pgxpool
func InitDatabase(connString string) error {
	var err error
	pool, err = pgxpool.New(context.Background(), connString)
	return err
}

// Close function to shutdown gracefully
func Close() {
	if pool != nil {
		pool.Close()
	}
}

// Insert cashflow values for 2020 - 2023 years
func InsertFCF(ticker string, cashFlow2020, cashFlow2021, cashFlow2022, cashFlow2023 string) error {
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

	_, err = pool.Exec(context.Background(), query, ticker, cf2020, cf2021, cf2022, cf2023)
	if err != nil {
		return fmt.Errorf("failed to insert data into database: %v", err)
	}

	return nil
}
