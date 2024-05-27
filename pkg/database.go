package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func InitDatabase(connString string) error {
	var err error
	pool, err = pgxpool.New(context.Background(), connString)
	return err
}

func InsertData(ticker string, cashFlow2020, cashFlow2021, cashFlow2022, cashFlow2023 int) error {
	query := `
        INSERT INTO stock_cash_flow 
            (ticker, cash_flow_2020, cash_flow_2021, cash_flow_2022, cash_flow_2023) 
        VALUES 
            ($1, $2, $3, $4, $5)`

	_, err := pool.Exec(context.Background(), query, ticker, cashFlow2020, cashFlow2021, cashFlow2022, cashFlow2023)
	return err
}
