package kstock

// Querier interface defines how to retrieve KOSPI and KOSDAQ stock list and price
type Querier interface {
	GetKospiStockList() ([]StockInfo, error)
	GetKosdaqStockList() ([]StockInfo, error)
	GetStockPrice(string, int) ([]StockData, error)
}
