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
	Pretax_Income    string `json:"Pretax_Income"`

	// Balance Sheet
	Total_Debt string `json:"Total_Debt"`

	// Summary Sheet
	Beta       string `json:"Beta"`
	Market_Cap string `json:"Market_Cap"`
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
	fmt.Println("Moving on to Requesting")
	fmt.Printf("Attempting to scrape data for ticker: %s\n", ticker)

	// Error Handle if not correct website ticker / other error
	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Request URL: %s failed with status code: %d, error: %v\n", r.Request.URL, r.StatusCode, err)
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

	c.OnHTML("div.tableBody", func(e *colly.HTMLElement) {
		newItem := newItem{
			Interest_Expense: cleanAndParseFCF(e.ChildText("div.row:nth-child(21) div:nth-child(3)")),
			Pretax_Income:    cleanAndParseFCF(e.ChildText("div.row:nth-child(8) div:nth-child(3)")),
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

	c.OnHTML("div.tableBody", func(e *colly.HTMLElement) {

		// Iterate over each row
		e.ForEach("div.row", func(i int, row *colly.HTMLElement) {

			// Check for the rowTitle with the desired title
			rowTitle := row.DOM.Find("div.rowTitle")
			titleAttr := rowTitle.AttrOr("title", "")

			if titleAttr == "Total Debt" {

				// Find the next column div and get its text
				columns := row.DOM.Find("div.column")
				if columns.Length() > 1 {
					nextColumn := columns.Eq(1) // Get the second column
					nextColumnText := nextColumn.Text()
					totalDebtValue := cleanAndParseFCF(nextColumnText)

					// Set value to strut
					newItem := newItem{
						Total_Debt: totalDebtValue,
					}
					items = append(items, newItem)
				} else {
					fmt.Println("Nothing in next row")
				}
			} else {
				fmt.Println("Not the target row")
			}
		})
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

func ScrapeSummary(ticker string) ([]newItem, error) {

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
		fmt.Println("// Summary // Visiting", r.URL.String())
	})
	// Error Handle if not correct website ticker / other error
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	// Confirm Response
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("// Summary // Visited", r.Request.URL)
	})
	// Scrape FCF, goes to last row of table body div, and selects chosen childs
	c.OnHTML("div.yf-mrt107", func(e *colly.HTMLElement) {
		newItem := newItem{
			Beta:       cleanAndParseFCF(e.ChildText("li.yf-mrt107:nth-child(10) > span:nth-child(2)")),
			Market_Cap: detectMarketCap((e.ChildText("li.yf-mrt107:nth-child(9) > span:nth-child(2) > fin-streamer:nth-child(1)"))),
		}
		items = append(items, newItem)
	})
	// Confirmed vist and done filling out OnHTML callback
	c.OnScraped(func(r *colly.Response) {
		fmt.Println("// Summary // Finished", r.Request.URL)
		fmt.Println()
	})
	// URL setup from User input
	url := fmt.Sprintf("https://finance.yahoo.com/quote/%s", ticker)
	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	return items, nil

}

// Function to clean and parse FCF strings
func cleanAndParseFCF(fcfString string) string {
	// Remove commas and trim spaces
	fcfString = strings.ReplaceAll(fcfString, ",", "")
	fcfString = strings.TrimSpace(fcfString)

	// Parse the cleaned string as a float
	parsedFCF, err := strconv.ParseFloat(fcfString, 64)
	if err != nil {
		fmt.Printf("Parsing error: %v\n", err)
		return ""
	}

	return fmt.Sprintf("%.2f", parsedFCF)
}

// Function to detect Market Cap is in trillions(T), billions(B), millions(M)

func detectMarketCap(marketCap string) string {

	trimmed := marketCap[:len(marketCap)-1]
	floatMarketCap, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		fmt.Printf("Parsing error: %v\n", err)
		return ""
	}
	if strings.Contains(marketCap, "T") {
		fmt.Println(floatMarketCap)
		return fmt.Sprintf("%.0f", floatMarketCap*1000000000000)
	}
	if strings.Contains(marketCap, "B") {
		fmt.Println(floatMarketCap)
		return fmt.Sprintf("%.0f", floatMarketCap*1000000000)
	}
	if strings.Contains(marketCap, "M") {
		fmt.Println(floatMarketCap)
		return fmt.Sprintf("%.0f", floatMarketCap*1000000)
	}
	return fmt.Sprintf("%.0f", floatMarketCap)
}
