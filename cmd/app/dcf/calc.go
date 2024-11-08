package calc

import (
	"context"
	"log"

	database "github.com/Tawxyn/goStockScraper/pkg"
)

type FinancialService struct {
	db *database.Postgres
}

func NewFinancialService(db *database.Postgres) *FinancialService {
	return &FinancialService{db: db}
}

func (fs *FinancialService) CalculateWAAC(ticker string) error {
	ctx := context.Background()

	financialData, err := fs.db.GetFinancials(ctx, ticker)
	if err != nil {
		log.Fatalf("Error retrieving finanical data: %v", err)
	}

	// Temporarily print all the financial data for debugging
	log.Printf("Financial Data for %s:\n", ticker)
	log.Printf("Cash Flow 2023: %f\n", financialData.CashFlow2023)
	log.Printf("Interest Expense: %f\n", financialData.InterestExpense2023)
	log.Printf("Total Debt: %f\n", financialData.TotalDebt2023)
	log.Printf("Pretax Income: %f\n", financialData.PretaxIncome2023)
	log.Printf("Beta: %f\n", financialData.BetaRecent)
	log.Printf("Market Cap: %f\n", financialData.MarketCap)

	return nil

}
