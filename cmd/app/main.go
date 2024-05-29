package main

import (
	"fmt"
	"log"
	"os"

	database "github.com/Tawxyn/goStockScraper/pkg"
	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
)

// Json item structure for scalability / orginization
type newItem struct {
	FCF_Year1 string `json:"FCF1"`
	FCF_Year2 string `json:"FCF2"`
	FCF_Year3 string `json:"FCF3"`
	FCF_Year4 string `json:"FCF4"`
}

func main() {

	//Load environment variables from .env file
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v\n", err)
	}

	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		log.Fatalf("DATABASE_URL was not found in the .env or is empty")
	}
	fmt.Println(connString)

	// Database import (pkg/database.go)
	err = database.InitDatabase(connString)
	if err != nil {
		log.Fatalf("Error in initializing database post .env load: %v\n", err)
	}
	defer database.Close() // Close database after main exists

	ticker := tickerInput()
	// Initiate new collector
	c := colly.NewCollector(
		// Whitelist website for visit
		colly.AllowedDomains("www.finance.yahoo.com", "finance.yahoo.com"),
	)
	// User agent to not get blocked
	// **TODO** randomize prceduraly Generate user agent to not be blocked in the future.
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"

	// Item sruct slice
	var items []newItem
	// Prior vist, request
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	// Error Handle if not correct website ticker / other error
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	// Confirm Response
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})
	// Scrape FCF, goes to last row of table body div, and selects chosen childs
	c.OnHTML("div.tableBody div.row:last-of-type", func(e *colly.HTMLElement) {
		newItem := newItem{
			FCF_Year1: e.ChildText("div:nth-child(3)"),
			FCF_Year2: e.ChildText("div:nth-child(4)"),
			FCF_Year3: e.ChildText("div:nth-child(5)"),
			FCF_Year4: e.ChildText("div:nth-child(6)"),
		}
		items = append(items, newItem)
	})
	// Confirmed vist and done filling out OnHTML callback
	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})
	// URL setup from User input
	url := fmt.Sprintf("https://finance.yahoo.com/quote/%s/cash-flow", ticker)
	c.Visit(url)
	fmt.Println(items)

	database.InsertFCF(ticker, "bruh", items[0].FCF_Year2, items[0].FCF_Year3, items[0].FCF_Year4)
	if err != nil {
		fmt.Printf("Error insertintg data into database: %v\n", err)
	} else {
		fmt.Println("Data insert succesffuly")
	}

}

// Obtain user ticker info.
func tickerInput() string {
	var input string
	fmt.Print("Input a stock ticker to analze: ")
	fmt.Scanln(&input)

	return input

}
