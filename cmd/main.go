package main

import (
	"fmt"
	stockQuerier "github.com/jinarusha/kstock/pkg/querier"
)

func main() {
	querier := stockQuerier.DaumQuerier{}

	perPage := 10
	page := 1

	quotes, err := querier.GetQuotes("A217270", perPage, page)
	fmt.Println(quotes, err)
}
