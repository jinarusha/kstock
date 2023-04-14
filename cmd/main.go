package main

import (
	"fmt"
	stockQuerier "github.com/jinarusha/kstock/pkg/querier"
)

func main() {
	querier := stockQuerier.DaumQuerier{}
	quotes, err := querier.GetQuotes("A217270", 1)
	fmt.Println(quotes, err)
}
