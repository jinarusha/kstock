package kstock

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const daumKospiPageUrl string = "http://finance.daum.net/quote/all.daum?type=S&stype=P"
const kospiMarket string = "KOSPI"

const daumKosdaqPageUrl string = "http://finance.daum.net/quote/all.daum?type=S&stype=Q"
const kosdaqMarket string = "KOSDAQ"

// DaumQuerier scrapes stock list and price from Daum portal
type DaumQuerier struct {
}

// KospiPageUrl fetches  KOSPI stock list page URL
func (q *DaumQuerier) KospiPageUrl() string {
	return daumKospiPageUrl
}

// KosdaqPageUrl fetches KOSDAQ stock list page URL
func (q *DaumQuerier) KosdaqPageUrl() string {
	return daumKosdaqPageUrl
}

// GetKospiStockList fetches KOSPI stock list
func (q *DaumQuerier) GetKospiStockList() ([]StockInfo, error) {
	body, err := getHtmlReader(q.KospiPageUrl())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	stockInfoList := []StockInfo{}

	numRows := doc.Find(".gTable tr").Length()
	doc.Find(".gTable tr").Each(func(i int, s *goquery.Selection) {
		// skip first 2 rows and last row
		if i < 2 || i >= numRows-1 {
			return
		}

		// KOSPI 종목
		s.Find("td").Each(func(j int, innerSel *goquery.Selection) {
			stockInfo := StockInfo{}
			if j%3 == 0 {
				code, exists := innerSel.Find("a").Attr("href")
				if !exists {
					return
				}
				codeArr := strings.SplitAfter(code, "=")
				stockInfo.Code = codeArr[1]
				stockInfo.Name = innerSel.Text()
				stockInfo.Market = kospiMarket
				stockInfoList = append(stockInfoList, stockInfo)
			}
		})
	})
	return stockInfoList, nil
}

// GetKosdaqStockList fetches KOSDAQ stock list
func (q *DaumQuerier) GetKosdaqStockList() ([]StockInfo, error) {
	body, err := getHtmlReader(q.KosdaqPageUrl())
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	stockInfoList := []StockInfo{}

	numRows := doc.Find(".gTable tr").Length()
	doc.Find(".gTable tr").Each(func(i int, s *goquery.Selection) {
		// skip first 2 rows and last row
		if i < 2 || i >= numRows-1 {
			return
		}

		// KOSDAQ 종목
		s.Find("td").Each(func(j int, innerSel *goquery.Selection) {
			if j%3 == 0 {
				code, exists := innerSel.Find("a").Attr("href")
				if !exists {
					return
				}
				stockInfo := StockInfo{}
				codeArr := strings.SplitAfter(code, "=")
				stockInfo.Code = codeArr[1]
				stockInfo.Name = innerSel.Text()
				stockInfo.Market = kosdaqMarket
				stockInfoList = append(stockInfoList, stockInfo)
			}
		})
	})
	return stockInfoList, nil
}

const daumStockDetailUrl string = "http://finance.daum.net/item/quote_yyyymmdd_sub.daum?code=%s&page=%d"

// StockDetailUrl fetches stock detail URL
func (q *DaumQuerier) StockDetailUrl(code string, page int) string {
	return fmt.Sprintf(daumStockDetailUrl, code, page)
}

// GetStockPrice fetches stock price data
func (q *DaumQuerier) GetStockPrice(code string, page int) ([]StockData, error) {

	stockDataList := []StockData{}

	if code == "" {
		return stockDataList, errors.New("code missing")
	}

	body, err := getHtmlReader(q.StockDetailUrl(code, page))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	// doc, err := goquery.NewDocument(fmt.Sprintf(DaumStockInfoListUrl, code, page))
	// if err != nil {
	// 	return stockDataList, err
	// }

	// empty page
	rowNum := doc.Find("#bbsList tr").Length()
	if rowNum == 1 {
		return stockDataList, errors.New("empty page")
	}

	doc.Find("#bbsList tr").Each(func(i int, rowSel *goquery.Selection) {
		// skip first two rows and last row
		if i < 2 || i == rowNum-1 {
			return
		}

		// skip empty rows
		cellLength := rowSel.Find("td").Length()
		if cellLength < 2 {
			return
		}

		stockData := StockData{
			Code: code,
		}

		rowSel.Find("td").Each(func(j int, innerSel *goquery.Selection) {
			// date
			if j == 0 {
				date, e := time.Parse("06.01.02", innerSel.Text())
				if e != nil {
					return
				}
				stockData.Date = date
			}
			// opening price
			if j == 1 {
				val, err := strconv.Atoi(strings.Replace(innerSel.Text(), ",", "", -1))
				if err != nil {
					return
				}
				stockData.OpenPrice = val
			}
			// high price
			if j == 2 {
				val, err := strconv.Atoi(strings.Replace(innerSel.Text(), ",", "", -1))
				if err != nil {
					return
				}
				stockData.HighPrice = val
			}
			// low price
			if j == 3 {
				val, err := strconv.Atoi(strings.Replace(innerSel.Text(), ",", "", -1))
				if err != nil {
					return
				}
				stockData.LowPrice = val
			}
			// closing price
			if j == 4 {
				val, err := strconv.Atoi(strings.Replace(innerSel.Text(), ",", "", -1))
				if err != nil {
					return
				}
				stockData.ClosingPrice = val
			}
			// volume
			if j == 7 {
				val, err := strconv.Atoi(strings.Replace(innerSel.Text(), ",", "", -1))
				if err != nil {
					return
				}
				stockData.Volume = val
			}
		})

		// if empty row, skip
		if stockData.Date.IsZero() {
			return
		}

		stockDataList = append(stockDataList, stockData)
	})

	return stockDataList, nil
}

func getHtmlReader(url string) (io.Reader, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(body), nil

}
