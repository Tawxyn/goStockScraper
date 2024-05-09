package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

type item struct {
	Free_Cash_Flow string
	Current_Price  string
	Current_Low    string
	Current_High   string
}

func main() {
	ticker := tickerInput()
	// Initiate new collector
	c := colly.NewCollector(
		// Whitelist website for visit
		colly.AllowedDomains("finance.yahoo.com"),
	)
	// Prior vist, request
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit("https://finance.yahoo.com/" + ticker)
}

// Obtain user ticker info.
func tickerInput() string {
	var input string
	fmt.Println("Input a stock ticker to analze: ")
	fmt.Scanln(&input)

	return input

}
