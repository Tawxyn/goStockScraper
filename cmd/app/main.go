package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gocolly/colly"
	"github.com/jackc/pgx/v5"
)

// Json item structure for scalability / orginization
type item struct {
	FCF_Year1 string `json:"FCF1"`
	FCF_Year2 string `json:"FCF2"`
	FCF_Year3 string `json:"FCF3"`
	FCF_Year4 string `json:"FCF4"`
}

func main() {
	// .env pull for pgx / DB connection
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))

	ticker := tickerInput()
	// Initiate new collector
	c := colly.NewCollector(
		// Whitelist website for visit
		colly.AllowedDomains("www.finance.yahoo.com", "finance.yahoo.com"),
	)
	// User agent to not get blocked
	// **TODO** randomize prceduraly Generate user agent to not be blocked in the future.
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"

	// Item Truct slice
	items := []item{}
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
		item := item{
			FCF_Year1: e.ChildText("div:nth-child(3)"),
			FCF_Year2: e.ChildText("div:nth-child(4)"),
			FCF_Year3: e.ChildText("div:nth-child(5)"),
			FCF_Year4: e.ChildText("div:nth-child(6)"),
		}
		items = append(items, item)
	})
	// Confirmed vist and done filling out OnHTML callback
	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})
	// URL setup from User input
	url := fmt.Sprintf("https://finance.yahoo.com/quote/%s/cash-flow", ticker)
	c.Visit(url)
	fmt.Println(items)
}

// Obtain user ticker info.
func tickerInput() string {
	var input string
	fmt.Print("Input a stock ticker to analze: ")
	fmt.Scanln(&input)

	return input

}
