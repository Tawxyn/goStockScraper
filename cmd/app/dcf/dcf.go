package dcf

import (
	database "github.com/Tawxyn/goStockScraper/pkg"
)

func calculateWAAC(ticker string) {
	database.GetFinancials(ticker)
}
