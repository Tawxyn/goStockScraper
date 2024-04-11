package main

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

func main() {
	// Initiate new collector
	c := colly.NewCollector(
		// Whitelist website for visit
		colly.AllowedDomains("finance.yahoo.com")
	)

}
