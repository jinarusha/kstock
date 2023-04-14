package entity

import "time"

// Quote stores entity price data for a particular day
type Quote struct {
	Code         string
	Date         string
	DateAsNum    time.Time
	OpenPrice    float64
	HighPrice    float64
	LowPrice     float64
	ClosingPrice float64
	Volume       int
}
