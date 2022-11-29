package stocks

func GetAllClosings(ticker string, data map[string][]StockRow) []float64 {
	closings := []float64{}
	// fmt.Println(len(data[ticker]))
	for _, row := range data[ticker] {
		// fmt.Println(row)
		closings = append(closings, row.Close)
	}
	return closings
}

func GetAllOpenings(ticker string, data map[string][]StockRow) []float64 {
	openings := []float64{}
	for _, row := range data[ticker] {
		openings = append(openings, row.Open)
	}
	return openings
}
