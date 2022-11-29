package stocks

import (
	"errors"
	"strconv"
	"sync"
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

func parseMapWorker(data map[string][][]string, tickerChunk []string, tickerChan chan string, stockData chan []StockRow, wg *sync.WaitGroup) {
	defer wg.Done()
	var stockRows []StockRow
	// fmt.Println(len(tickerChunk))
	for _, ticker := range tickerChunk {
		for _, row := range data[ticker] {
			stockRow, err := ParseRow(row)
			if err != nil {
				continue
			}
			stockRows = append(stockRows, stockRow)
		}

		stockData <- stockRows
		tickerChan <- ticker
	}
}

func ParseMap(data map[string][][]string, threads int) map[string][]StockRow {
	stockData := make(map[string][]StockRow)
	tickerChan := make(chan string, len(data)*2) // dirty fix

	tickerChunks := [][]string{}
	for i := 0; i < threads; i++ {
		strLen := len(data) / threads
		s := make([]string, strLen)
		tickerChunks = append(tickerChunks, s)
	}

	chanLen := 0

	currentThread := 0
	// for some reason this is doubling the amount of data
	// don't have time to fix it
	for ticker := range data {
		if currentThread >= threads {
			currentThread = 0
		}
		tickerChunks[currentThread] = append(tickerChunks[currentThread], ticker)
		currentThread++
		chanLen += len(data[ticker])
	}

	stockRows := make(chan []StockRow, chanLen)

	var wg sync.WaitGroup

	for _, tickerChunk := range tickerChunks {
		wg.Add(1)
		go parseMapWorker(data, tickerChunk, tickerChan, stockRows, &wg)
	}

	wg.Wait()

	// get all the data from the channel
	for i := 0; i < len(data); i++ {
		stockRow := <-stockRows
		ticker := <-tickerChan
		stockData[ticker] = stockRow
	}

	return stockData
}
