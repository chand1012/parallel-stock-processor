package stocks

import (
	"errors"
	"strconv"
	"time"
)

type StockRow struct {
	Date   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
	Ticker string
}

const dateFormat = "2006-01-02"

func StrSliceContains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

func ParseRow(row []string) (StockRow, error) {
	var stockRow StockRow

	if len(row[1]) < 10 {
		return stockRow, errors.New("invalid date")
	}

	stockRow.Ticker = row[0]

	date, err := time.Parse(dateFormat, row[1][0:10])
	if err != nil {
		return stockRow, err
	}

	stockRow.Date = date

	stockRow.Open, err = strconv.ParseFloat(row[2], 64)
	if err != nil {
		return stockRow, err
	}

	stockRow.High, err = strconv.ParseFloat(row[3], 64)
	if err != nil {
		return stockRow, err
	}

	stockRow.Low, err = strconv.ParseFloat(row[4], 64)
	if err != nil {
		return stockRow, err
	}

	stockRow.Close, err = strconv.ParseFloat(row[5], 64)
	if err != nil {
		return stockRow, err
	}

	stockRow.Volume, err = strconv.ParseFloat(row[6], 64)
	if err != nil {
		return stockRow, err
	}

	return stockRow, nil
}
