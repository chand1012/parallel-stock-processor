package stocks

func ParseCSV(data [][]string) ([]StockRow, error) {
	var stockRows []StockRow

	for i, row := range data {
		// skip the first row
		if i == 0 {
			continue
		}
		stockRow, err := ParseRow(row)
		if err != nil {
			return nil, err
		}

		stockRows = append(stockRows, stockRow)
	}

	return stockRows, nil
}
