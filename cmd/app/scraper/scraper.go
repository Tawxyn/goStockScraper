package scraper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

// Json item structure for scalability / orginization
type newItem struct {
	// Cash Flows
	FCF_Year1 string `json:"FCF1"`
	FCF_Year2 string `json:"FCF2"`
	FCF_Year3 string `json:"FCF3"`
	FCF_Year4 string `json:"FCF4"`

	// Income Statements
	Interest_Expense string `json:"Interest_Expense"`

	// Balance Sheet
	Total_Debt string `json:"Total_Debt"`
}

func ScrapeCashFlow(ticker string) ([]newItem, error) {

	// Continue with your application logic here
	fmt.Println("Continuing with application logic...")

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
		fmt.Println("// Cash Flow // Visiting", r.URL.String())
	})
	// Error Handle if not correct website ticker / other error
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	// Confirm Response
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("// Cash Flow // Visited", r.Request.URL)
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
		fmt.Println("// Cash Flow // Finished", r.Request.URL)
		fmt.Println()
	})
	// URL setup from User input
	url := fmt.Sprintf("https://finance.yahoo.com/quote/%s/cash-flow", ticker)
	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	return items, nil

}

func ScrapeIncomeStatement(ticker string) ([]newItem, error) {

	c := colly.NewCollector(

		colly.AllowedDomains("www.finance.yahoo.com", "finance.yahoo.com"),
	)

	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"

	var items []newItem

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("// Income Statement // Visiting", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("// Income Statement // Visited", r.Request.URL)
	})

	c.OnHTML("div.tableBody div.row:nth-child(21)", func(e *colly.HTMLElement) {
		newItem := newItem{
			Interest_Expense: cleanAndParseFCF(e.ChildText("div:nth-child(3)")),
		}
		items = append(items, newItem)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("// Income Statement // Finished", r.Request.URL)
		fmt.Println()
	})
	// URL setup from User input
	url := fmt.Sprintf("https://finance.yahoo.com/quote/%s/financials/", ticker)
	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	return items, nil

}

func ScrapeBalanceSheet(ticker string) ([]newItem, error) {

	c := colly.NewCollector(

		colly.AllowedDomains("www.finance.yahoo.com", "finance.yahoo.com"),
	)

	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"

	var items []newItem

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("// Balance Sheet // Visiting ", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("// Balance Sheet // Visited", r.Request.URL)
	})

	c.OnHTML("div.tableBody div.row:nth-child(10)", func(e *colly.HTMLElement) {
		newItem := newItem{
			Total_Debt: cleanAndParseFCF(e.ChildText("div:nth-child(2)")),
		}
		items = append(items, newItem)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("// Balance Sheet // Finished", r.Request.URL)
		fmt.Println()
	})
	// URL setup from User input
	url := fmt.Sprintf("https://finance.yahoo.com/quote/%s/balance-sheet/", ticker)
	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	return items, nil

}

// Function to clean and parse FCF strings
func cleanAndParseFCF(fcfString string) string {
	// Remove commas from the string
	fcfString = strings.ReplaceAll(fcfString, ",", "")

	// Attempt to parse the cleaned string as a float
	parsedFCF, err := strconv.ParseFloat(fcfString, 64)
	if err != nil {
		// If parsing fails, return an empty string
		return ""
	}

	// Convert the float to a string (without decimal places)
	return fmt.Sprintf("%.0f", parsedFCF)
}
