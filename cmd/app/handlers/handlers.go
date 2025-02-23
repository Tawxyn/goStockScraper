package handlers

import (
	"context"
	"fmt"
	"html/template"
	"net/http"

	calc "github.com/Tawxyn/goStockScraper/cmd/app/dcf"
	scraper "github.com/Tawxyn/goStockScraper/cmd/app/scraper"
	users "github.com/Tawxyn/goStockScraper/cmd/users"
	database "github.com/Tawxyn/goStockScraper/pkg"
)

type Handler struct {
	db               *database.Postgres
	FinancialService *calc.FinancialService
}

func NewHandler(db *database.Postgres, financialService *calc.FinancialService) *Handler {
	return &Handler{
		db:               db,
		FinancialService: financialService,
	}
}

func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "views/index.html")
}

func (h *Handler) AnalyzeHandler(w http.ResponseWriter, r *http.Request) {
	ticker := r.URL.Query().Get("stockSymbol")

	if ticker == "" {
		http.Error(w, "Missing stockSymbol parameter", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	// Check if the ticker exists in the database
	exists, err := h.db.CheckTickerExists(ctx, ticker)
	if err != nil {
		fmt.Println("Database error in CheckTickerExists:", err) // Debugging
		http.Error(w, "Error checking ticker", http.StatusInternalServerError)
		return
	}

	if exists {
		fmt.Fprintf(w, "Ticker %s already exists.", ticker)
		return
	}

	// Scrape finanical data
	cashFlowitems, err := scraper.ScrapeCashFlow(ticker)
	if err != nil {
		fmt.Println("Error scraping Cash Flow:", err) // Debugging
		http.Error(w, "Error scraping Cash Flow Page", http.StatusInternalServerError)
		return
	}

	incomeStatementItems, err := scraper.ScrapeIncomeStatement(ticker)
	if err != nil {
		fmt.Println("Error scraping Income Statement:", err) // Debugging
		http.Error(w, "Error scraping Income Statement", http.StatusInternalServerError)
		return
	}

	totalDebtItems, err := scraper.ScrapeBalanceSheet(ticker)
	if err != nil {
		fmt.Println("Error scraping Balance Sheet:", err) // Debugging
		http.Error(w, "Error scraping Balance Sheet", http.StatusInternalServerError)
		return
	}

	summaryitems, err := scraper.ScrapeSummary(ticker)
	if err != nil {
		fmt.Println("Error scraping Summary:", err) // Debugging
		http.Error(w, "Error scraping Cash Flow Page", http.StatusInternalServerError)
		return
	}

	if len(cashFlowitems) > 0 && len(incomeStatementItems) > 0 && len(totalDebtItems) > 0 && len(summaryitems) > 0 {
		err = h.db.InsertFCF(ctx, ticker, cashFlowitems[0].FCF_Year1, cashFlowitems[0].FCF_Year2, cashFlowitems[0].FCF_Year3,
			cashFlowitems[0].FCF_Year4, incomeStatementItems[0].Interest_Expense, totalDebtItems[0].Total_Debt, incomeStatementItems[0].Pretax_Income,
			summaryitems[0].Beta, summaryitems[0].Market_Cap)
		if err != nil {
			fmt.Println("Database Insert Error:", err) // Debugging
			http.Error(w, "Error inserting data", http.StatusInternalServerError)
		} else {
			fmt.Fprintln(w, "Data successfully inserted")
		}
	} else {
		fmt.Fprintln(w, "No data found to insert")
	}
}

func (h *Handler) CalculateWAAC(w http.ResponseWriter, r *http.Request) {
	ticker := r.URL.Query().Get("stockSymbol")

	err := h.FinancialService.CalculateWAAC(ticker)
	if err != nil {

		http.Error(w, "Error with calculating WAAC", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Financial data for %s has been printed to the logs.", ticker)
}

// User Handler logic with templating

func (h *Handler) UserHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "views/login.html")
}

func getSignInPage(w http.ResponseWriter, r *http.Request) {
	templating(w, "sign-in.html", nil)
}

func getSignUpPage(w http.ResponseWriter, r *http.Request) {
	templating(w, "sign-up.html", nil)
}

func templating(w http.ResponseWriter, fileName string, data interface{}) {
	t, _ := template.ParseFiles(fileName)
	t.ExecuteTemplate(w, fileName, data)
}

func signInUser(w http.ResponseWriter, r *http.Request) {
	newUser := getUser(r)
	users.DefaultUserService.CreateUser(newUser)
}

func signUpUser(w http.ResponseWriter, r *http.Request) {

}

func getUser(r *http.Request) users.User {
	email := r.FormValue("email")
	password := r.FormValue("password")
	return users.User{
		Email:    email,
		Password: password,
	}
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/sign-in":
		signInUser(w, r)
	case "/sign-up":
		signUpUser(w, r)
	case "sign-in-form":
		getSignInPage(w, r)
	case "sign-up-form":
		getSignUpPage(w, r)
	}
}
