# kstock
[![GoDoc](https://godoc.org/github.com/jinarusha/kstock?status.svg)](https://godoc.org/github.com/jinarusha/kstock)
[![Go Report Card](https://goreportcard.com/badge/jinarusha/kstock)](http://goreportcard.com/report/jinarusha/kstock)
[![Coverage Status](http://img.shields.io/coveralls/jinarusha/kstock.svg)](https://coveralls.io/r/jinarusha/kstock)
[![Build Status](https://travis-ci.org/jinarusha/kstock.svg?branch=master)](https://travis-ci.org/jinarusha/kstock)

kstock is a simple stock price scraper library that parses KOSPI and KOSDAQ pages on major Korean portals. Only scraping from Daum at the moment.
Kstock는 주요 포털에서 KOSPI 그리고 KOSDAQ 종목 페이지를 파싱하여 가격 정보를 가져오는 라이브러리입니다.

## Portal List
- Daum / 다음

## Installation

```
$ go get github.com/jinarusha/kstock
```

## Usage

```
import "github.com/jinarusha/kstock"

func main() {
  querier := kstock.DaumQuerier{}
  stockDataList, err := querier.GetStockPrice("017960", 1)
  fmt.Println(stockDataList, err)
}
```

## Test

```
$ go get github.com/jinarusha/kstock
```

### To see test results in browser using goconvey

```
$ goconvey
```

All pull requests and bug reports are welcome.

## License

Released under the [MIT License](https://github.com/jinarusha/kstock/blob/master/LICENSE).
