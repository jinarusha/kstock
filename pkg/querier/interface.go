package querier

import (
	"github.com/jinarusha/kstock/pkg/entity"
)

// Querier interface defines how to retrieve KOSPI and KOSDAQ entity list and price
type Querier interface {
	GetStocks(entity.MarketType) ([]entity.Stock, error)
	GetQuotes(string, int, int) ([]entity.Quote, error)
}
