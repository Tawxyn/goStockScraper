package main

import (
	"fmt"
	"log"

	"github.com/gocolly/colly"
)

type Item struct {
	FCF string 'json:"FCF"'
}

func main() {
	//ticker := tickerInput()
	// Initiate new collector
	c := colly.NewCollector(
		// Whitelist website for visit
		colly.AllowedDomains("www.finance.yahoo.com", "finance.yahoo.com"),
	)
	// Prior vist, request
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	// Error Handle if not correct website ticker / other error
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Error while scraping %s: %v\n", r.Request.URL, err)
	})
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})
	// Scrape FCF
	c.OnHTML("main div[column.svelte-1xjz32c]", func(e *colly.HTMLElement) {
		fmt.Println("Free Cash Flow = ", e.Text)
	})

	//url := fmt.Sprintf("https://finance.yahoo.com/quote/%s/cash-flow", ticker)
	c.Visit("https://finance.yahoo.com/quote/INTC/cash-flow")

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})
}

// Obtain user ticker info.
func tickerInput() string {
	var input string
	fmt.Println("Input a stock ticker to analze: ")
	fmt.Scanln(&input)

	return input

}
