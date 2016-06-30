package kstock

import "time"

// StockInfo stores general stock info
type StockInfo struct {
	Code   string `bson:"code" json:"code"`
	Name   string `bson:"n" json:"name"`
	Market string `bson:"m" json:"market"`
}

// StockData stores stock price data for a particular day
type StockData struct {
	Code         string    `bson:"code" json:"code"`
	Date         time.Time `bson:"date" json:"date"`
	OpenPrice    int       `bson:"o" json:"openPrice"`
	HighPrice    int       `bson:"h" json:"highPrice"`
	LowPrice     int       `bson:"l" json:"lowPrice"`
	ClosingPrice int       `bson:"c" json:"closingPrice"`
	Volume       int       `bson:"v" json:"volume"`
}
