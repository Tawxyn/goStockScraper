package handlers

import (
	"context"
	"fmt"
	"net/http"

	scraper "github.com/Tawxyn/goStockScraper/cmd/app/scraper"
	database "github.com/Tawxyn/goStockScraper/pkg"
)

type Handler struct {
	db *database.Postgres
}

func NewHandler(db *database.Postgres) *Handler {
	return &Handler{db: db}
}

func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "views/index.html")
}

func (h *Handler) AnalyzeHandler(w http.ResponseWriter, r *http.Request) {
	ticker := r.URL.Query().Get("stockSymbol")

	ctx := context.Background()

	exists, err := h.db.CheckTickerExists(ctx, ticker)
	if err != nil {
		http.Error(w, "Error checking ticker", http.StatusInternalServerError)
		return
	}

	if exists {
		fmt.Fprintf(w, "Ticker %s already exists.", ticker)
		return
	}

	// Handle Cash Flow
	cashFlowitems, err := scraper.ScrapeCashFlow(ticker)
	if err != nil {

		http.Error(w, "Error scraping Cash Flow Page", http.StatusInternalServerError)
		return
	}

	incomeStatementItems, err := scraper.ScrapeIncomeStatement(ticker)
	if err != nil {

		http.Error(w, "Error scraping Income Statement", http.StatusInternalServerError)
		return
	}

	totalDebtItems, err := scraper.ScrapeBalanceSheet(ticker)
	if err != nil {

		http.Error(w, "Error scraping Balance Sheet", http.StatusInternalServerError)
		return
	}

	if len(cashFlowitems) > 0 && len(incomeStatementItems) > 0 && len(totalDebtItems) > 0 {
		err = h.db.InsertFCF(ctx, ticker, cashFlowitems[0].FCF_Year1, cashFlowitems[0].FCF_Year2, cashFlowitems[0].FCF_Year3,
			cashFlowitems[0].FCF_Year4, incomeStatementItems[0].Interest_Expense, totalDebtItems[0].Total_Debt)
		if err != nil {
			http.Error(w, "Error inserting data", http.StatusInternalServerError)
		} else {
			fmt.Fprintln(w, "Data successfully inserted")
		}
	} else {
		fmt.Fprintln(w, "No data found to insert")
	}
}
