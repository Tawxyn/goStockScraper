package main

import (
	"github.com/gocolly/colly"
)

type item struct {
	Free_Cash_Flow string
	Current_Price  string
	Current_Low    string
	Current_High   string
}

func main() {

	url := "https://finance.yahoo.com/quote/TSLA/cash-flow"

	// Initiate new collector
	c := colly.NewCollector(
		// Whitelist website for visit
		colly.AllowedDomains("finance.yahoo.com"),
	)
}
