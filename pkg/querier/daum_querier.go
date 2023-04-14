package querier

import (
	"fmt"
	"github.com/jinarusha/kstock/pkg/entity"
	"github.com/jinarusha/kstock/pkg/httpclient"
	jsoniter "github.com/json-iterator/go"
	"io"
	"net/http"
	"time"
)

const (
	daumStockListPageUrlFormat string = "https://finance.daum.net/api/quotes/sectors?fieldName=&order=&perPage=&market=%s&page="
	daumQuoteURLFormat         string = "https://finance.daum.net/api/quote/%s/days?symbolCode=%s&page=%d&perPage=%d&pagination=true"
	daumStockDetailPerPage            = 40

	quoteDateFormat = "2006-01-02"
)

// DaumQuerier scrapes entity list and price from Daum portal
type DaumQuerier struct {
}

type daumStockListResp struct {
	Data []struct {
		IncludedStocks []struct {
			Name       string  `json:"name"`
			Code       string  `json:"code"`
			SymbolCode string  `json:"symbolCode"`
			TradePrice float64 `json:"tradePrice"`
		} `json:"includedStocks"`
	} `json:"data"`
}

func (q *DaumQuerier) GetStocks(marketType entity.MarketType) ([]entity.Stock, error) {
	return q.getStockList(marketType)
}

func (q *DaumQuerier) getStockList(marketType entity.MarketType) ([]entity.Stock, error) {
	headers := http.Header{}
	headers.Set("authority", "finance.daum.net")
	headers.Set("accept", "application/json")
	headers.Set("referer", "https://finance.daum.net/domestic/all_stocks?market="+string(marketType))
	headers.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")

	url := fmt.Sprintf(daumStockListPageUrlFormat, marketType)
	resp, err := httpclient.GetClient().Get(url, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var listResp daumStockListResp
	err = jsoniter.Unmarshal(respBody, &listResp)
	if err != nil {
		return nil, err
	}

	var stocks []entity.Stock
	for _, entry := range listResp.Data {
		for _, stockResp := range entry.IncludedStocks {
			stocks = append(stocks, entity.Stock{
				Code:   stockResp.SymbolCode,
				Name:   stockResp.Name,
				Market: marketType,
			})

		}
	}

	return stocks, nil
}

type daumPaginatedQuoteResp struct {
	Data []struct {
		SymbolCode        string  `json:"symbolCode"`
		Date              string  `json:"date"`
		TradePrice        float64 `json:"tradePrice"`
		OpeningPrice      float64 `json:"openingPrice"`
		HighPrice         float64 `json:"highPrice"`
		LowPrice          float64 `json:"lowPrice"`
		PeriodTradeVolume float64 `json:"periodTradeVolume"`
	} `json:"data"`
}

func (q *DaumQuerier) GetQuotes(symbolCode string, page int) ([]entity.Quote, error) {
	headers := http.Header{}
	headers.Set("authority", "finance.daum.net")
	headers.Set("accept", "application/json")
	headers.Set("referer", "https://finance.daum.net/quotes/"+symbolCode)
	headers.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")

	url := fmt.Sprintf(daumQuoteURLFormat, symbolCode, symbolCode, page, daumStockDetailPerPage)
	resp, err := httpclient.GetClient().Get(url, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var paginationResp daumPaginatedQuoteResp
	err = jsoniter.Unmarshal(respBody, &paginationResp)
	if err != nil {
		return nil, err
	}

	var quotes []entity.Quote
	for _, entry := range paginationResp.Data {
		date, err := time.Parse("2006-01-02 15:04:05", entry.Date)
		if err != nil {
			return nil, err
		}

		quotes = append(quotes, entity.Quote{
			Code:         entry.SymbolCode,
			Date:         date.Format(quoteDateFormat),
			DateAsNum:    date,
			OpenPrice:    entry.OpeningPrice,
			HighPrice:    entry.HighPrice,
			LowPrice:     entry.LowPrice,
			ClosingPrice: entry.TradePrice,
			Volume:       int(entry.PeriodTradeVolume),
		})
	}

	return quotes, nil
}
