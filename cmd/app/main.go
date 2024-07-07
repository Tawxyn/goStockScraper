package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"strconv"

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

	// Create a context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Ensure that the context is canceled when main returns

	//Load environment variables from .env file
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v\n", err)
	}

	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		log.Fatalf("DATABASE_URL was not found in the .env or is empty")
	}

	// Database import (pkg/database.go)
	pgInstance, err := database.InitDatabase(ctx, connString)
	if err != nil {
		log.Fatalf("Error initializing database post .env load: %v\n", err)
	}
	defer pgInstance.Close() // Close database after main exits

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
			FCF_Year4: cleanAndParseFCF(e.ChildText("div:nth-child(3)")),
			FCF_Year3: cleanAndParseFCF(e.ChildText("div:nth-child(4)")),
			FCF_Year2: cleanAndParseFCF(e.ChildText("div:nth-child(5)")),
			FCF_Year1: cleanAndParseFCF(e.ChildText("div:nth-child(6)")),
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

	// pgInstance to call InsertFCF and pass the context
	if len(items) > 0 {
		err = pgInstance.InsertFCF(ctx, ticker, items[0].FCF_Year1, items[0].FCF_Year2, items[0].FCF_Year3, items[0].FCF_Year4)
		if err != nil {
			fmt.Printf("Error inserting data into database: %v\n", err)
		} else {
			fmt.Println("Data successfully inserted")
		}
	} else {
		fmt.Println("No data found to insert")
	}

}

// Obtain user ticker info.
func tickerInput() string {
	var input string
	fmt.Print("Input a stock ticker to analze: ")
	fmt.Scanln(&input)

	return input

}

// Function to clean and parse FCF strings
func cleanAndParseFCF(fcfString string) string {
	// Remove commas from the string
	fcfString = strings.ReplaceAll(fcfString, ",", "")

	// Attempt to parse the cleaned string as a float
	parsedFCF, err := strconv.ParseFloat(fcfString, 64)
	if err != nil {
		// If parsing fails, return an empty string or handle the error as needed
		return ""
	}

	// Convert the float to a string (without decimal places)
	return fmt.Sprintf("%.0f", parsedFCF)
}

func DCFCalc() int {
	var input string
	fmt.Print("Would you like to run a DCF Analysis?")
	fmt.Scanln(&input)

	if input == "no" {

	} else {

	}
	return 0
}
