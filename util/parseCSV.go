package util

import (
	"github.com/chand1012/parallel-term-project/stocks"
)

func ParseCSV(data [][]string) ([]stocks.StockRow, error) {
	var stockRows []stocks.StockRow

	for i, row := range data {
		// skip the first row
		if i == 0 {
			continue
		}
		stockRow, err := stocks.ParseRow(row)
		if err != nil {
			return nil, err
		}

		stockRows = append(stockRows, stockRow)
	}

	return stockRows, nil
}
