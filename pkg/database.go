package database

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	db *pgxpool.Pool
}

// StockCashFlow struct
type StockData struct {
	// Cash Flow
	Ticker       string
	CashFlow2020 float64
	CashFlow2021 float64
	CashFlow2022 float64
	CashFlow2023 float64

	// Income Statement
	InterestExpense2023 float64
	PretaxIncome2023    float64

	// Balance Sheet
	TotalDebt2023 float64

	// Summary Sheet
	BetaRecent float64
	MarketCap  float64
}

var (
	pgInstance *Postgres
	pgOnce     sync.Once
)

// Initalize the database with pgxpool
func InitDatabase(ctx context.Context, connString string) (*Postgres, error) {
	var err error

	pgOnce.Do(func() {
		var db *pgxpool.Pool
		db, err = pgxpool.New(ctx, connString)
		if err == nil {
			pgInstance = &Postgres{db}
		}
	})

	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	return pgInstance, nil
}

func (pg *Postgres) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

// Close function to shutdown gracefully
func (pg *Postgres) Close() {
	pg.db.Close()
}

// InsertFCF inserts cash flow values for 2020-2023 years
func (pg *Postgres) InsertFCF(ctx context.Context, ticker string, cashFlow2020, cashFlow2021, cashFlow2022, cashFlow2023, InterestExpense2023, TotalDebt2023,
	PretaxIncome2023, BetaRecent, MarketCap string) error {
	var err error

	// Convert string values to integers
	cf2020, err := strconv.ParseFloat(cashFlow2020, 64)
	if err != nil {
		return fmt.Errorf("failed to convert cashFlow2020 to integer: %v", err)
	}

	cf2021, err := strconv.ParseFloat(cashFlow2021, 64)
	if err != nil {
		return fmt.Errorf("failed to convert cashFlow2021 to integer: %v", err)
	}

	cf2022, err := strconv.ParseFloat(cashFlow2022, 64)
	if err != nil {
		return fmt.Errorf("failed to convert cashFlow2022 to integer: %v", err)
	}

	cf2023, err := strconv.ParseFloat(cashFlow2023, 64)
	if err != nil {
		return fmt.Errorf("failed to convert cashFlow2023 to integer: %v", err)
	}

	interestexpense, err := strconv.ParseFloat(InterestExpense2023, 64)
	if err != nil {
		return fmt.Errorf("failed to convert InterestExpense2023 to integer: %v", err)
	}

	totaldebt, err := strconv.ParseFloat(TotalDebt2023, 64)
	if err != nil {
		return fmt.Errorf("failed to convert TotalDebt2023 to integer: %v", err)
	}

	pretaxincome, err := strconv.ParseFloat(PretaxIncome2023, 64)
	if err != nil {
		return fmt.Errorf("failed to convert PretaxIncome2023 to integer: %v", err)
	}

	beta, err := strconv.ParseFloat(BetaRecent, 64)
	if err != nil {
		return fmt.Errorf("failed to convert Beta to integer: %v", err)
	}

	marketcap, err := strconv.ParseFloat(MarketCap, 64)
	if err != nil {
		return fmt.Errorf("failed to convert Beta to integer: %v", err)
	}

	query := `
        INSERT INTO stock_info 
            (ticker, cash_flow_2020, cash_flow_2021, cash_flow_2022, cash_flow_2023, interest_expense, total_debt, pretax_income, beta, market_cap) 
        VALUES 
            ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err = pg.db.Exec(ctx, query, ticker, cf2020, cf2021, cf2022, cf2023, interestexpense, totaldebt, pretaxincome, beta, marketcap)
	if err != nil {
		return fmt.Errorf("failed to insert data into database: %v", err)
	}

	return nil
}

func (pg *Postgres) CheckTickerExists(ctx context.Context, ticker string) (bool, error) {
	var count int
	query := `SELECT COUNT(*)
			  FROM stock_info
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
	if count > 0 {
		fmt.Println("Code found in database")
	}
	return count > 0, nil
}

func (pg *Postgres) GetFinancials(ctx context.Context, ticker string) (*StockData, error) {

	financialData := &StockData{}

	query := `SELECT *
			  FROM stock_info
			  WHERE ticker = $1`

	err := pg.db.QueryRow(ctx, query, ticker).Scan(
		&financialData.Ticker,
		&financialData.CashFlow2020,
		&financialData.CashFlow2021,
		&financialData.CashFlow2022,
		&financialData.CashFlow2023,
		&financialData.InterestExpense2023,
		&financialData.PretaxIncome2023,
		&financialData.TotalDebt2023,
		&financialData.BetaRecent,
		&financialData.MarketCap)
	if err != nil {
		return nil, fmt.Errorf("error querying row for exporting financials: %w", err)
	}

	return financialData, nil
}
