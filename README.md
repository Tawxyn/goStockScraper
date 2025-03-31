# GoStockScraper

A high-performance financial data scraping and analysis tool built with Go that extracts, processes, and analyzes stock data from Yahoo Finance.

## Features

- **Data Collection**
  - Scrapes cash flow statements, income statements, balance sheets, and company summaries
  - Handles data normalization and cleaning
  - Implements anti-scraping measures including user agent rotation

- **Financial Analysis**
  - Discounted Cash Flow (DCF) valuation calculations
  - Weighted Average Cost of Capital (WACC) analysis
  - Market capitalization normalization (T/B/M to numeric)

- **Data Management**
  - PostgreSQL storage with connection pooling
  - Batch insert/update operations
  - Environment-based configuration

- **Web Interface**
  - REST API endpoints for financial analysis
  - Basic user authentication system
  - Static file serving for frontend integration

## Technical Stack

- **Backend**: Go 1.20+
- **Database**: PostgreSQL 15+
- **Scraping**: Colly framework
- **Database Driver**: PGX
- **Web Server**: Native net/http

## Installation

### Prerequisites

- Go 1.20+
- PostgreSQL 15+

```bash
# Clone the repository
git clone https://github.com/Tawxyn/goStockScraper.git
cd goStockScraper

# Install dependencies
go mod download

# Configure environment
cp .env.example .env
# Edit .env with your database credentials
```
## Database Setup
Create a PostgreSQL database

Update the .env file: 
DATABASE_URL="postgres://username:password@localhost:5432/dbname?sslmode=disable"

## Running the Application
```
# Development mode
go run cmd/main.go

# Production build
go build -o bin/goStockScraper cmd/main.go
./bin/goStockScraper
```
API Documentation
Endpoints
Endpoint	Method	Description	Parameters
/analyze	GET	Full financial analysis	stockSymbol (required)
/CalculateWAAC	POST	WACC calculation	JSON payload
/user	GET	User management endpoint	-
