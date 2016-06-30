package kstock

import (
	"testing"

	"github.com/jarcoal/httpmock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDaumQuerierStockList(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	daumQurier := DaumQuerier{}

	httpmock.RegisterResponder("GET", daumQurier.KospiPageUrl(),
		httpmock.NewStringResponder(200, readTestResponseFile("daum_kospi.html")))
	httpmock.RegisterResponder("GET", daumQurier.KosdaqPageUrl(),
		httpmock.NewStringResponder(200, readTestResponseFile("daum_kosdaq.html")))

	Convey("Parse KOSPI stock list page", t, func() {
		stockList, err := daumQurier.GetKospiStockList()
		So(err, ShouldEqual, nil)
		So(len(stockList), ShouldBeGreaterThan, 0)
	})

	Convey("Parse KOSDAQ stock list page", t, func() {
		stockList, err := daumQurier.GetKosdaqStockList()
		So(err, ShouldEqual, nil)
		So(len(stockList), ShouldBeGreaterThan, 0)
	})
}

func TestDaumQuerierStockDataList(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	daumQurier := DaumQuerier{}

	samsungElectronicsCode := "005930"
	page := 1

	httpmock.RegisterResponder("GET", daumQurier.StockDetailUrl(samsungElectronicsCode, page),
		httpmock.NewStringResponder(200, readTestResponseFile("daum_kospi_samsung_electronics.html")))

	Convey("Parse Samsung Electronics stock detail", t, func() {
		stockDataList, err := daumQurier.GetStockPrice(samsungElectronicsCode, page)
		So(err, ShouldEqual, nil)
		So(len(stockDataList), ShouldEqual, 30)

		latestStock := stockDataList[0]
		So(latestStock.Code, ShouldEqual, samsungElectronicsCode)
		So(latestStock.OpenPrice, ShouldBeGreaterThan, 0)
		So(latestStock.ClosingPrice, ShouldBeGreaterThan, 0)
		So(latestStock.HighPrice, ShouldBeGreaterThan, 0)
		So(latestStock.LowPrice, ShouldBeGreaterThan, 0)
		So(latestStock.Volume, ShouldBeGreaterThan, 0)

	})

}
