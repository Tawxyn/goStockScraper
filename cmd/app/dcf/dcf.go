package dcf

import (
	"context"
	"log"

	database "github.com/Tawxyn/goStockScraper/pkg"
)

type FinanicalService struct {
	db *database.Postgres
}

func (fs *FinanicalService) calculateWAAC(ticker string) {
	ctx := context.Background()

	financialData, err := fs.db.GetFinancials(ctx, ticker)
	if err != nil {
		log.Fatalf("Error retrieving finanical data: %v", err)
	}

	log.Printf("Ticker: %s, Cash Flow 2023: %f\n", financialData.Ticker, financialData.CashFlow2023)
}
