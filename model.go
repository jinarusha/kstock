package kstock

import "time"

// StockInfo stores general stock info
type StockInfo struct {
	Code   string
	Name   string
	Market string
}

// StockData stores stock price data for a particular day
type StockData struct {
	Code         string
	Date         time.Time
	OpenPrice    int
	HighPrice    int
	LowPrice     int
	ClosingPrice int
	Volume       int
}
