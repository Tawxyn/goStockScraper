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
	http.ServeFile(w, r, "../../views/index.html")
}

func (h *Handler) AnalyzeHandler(w http.ResponseWriter, r *http.Request) {
	ticker := r.URL.Query().Get("symbol")
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

	items, err := scraper.ScrapeFCF(ticker)
	if err != nil {
		http.Error(w, "Error scraping data", http.StatusInternalServerError)
		return
	}

	if len(items) > 0 {
		err = h.db.InsertFCF(ctx, ticker, items[0].FCF_Year1, items[0].FCF_Year2, items[0].FCF_Year3, items[0].FCF_Year4)
		if err != nil {
			http.Error(w, "Error inserting data", http.StatusInternalServerError)
		} else {
			fmt.Fprintln(w, "Data successfully inserted")
		}
	} else {
		fmt.Fprintln(w, "No data found to insert")
	}
}
