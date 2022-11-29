package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/montanaflynn/stats"

	"github.com/chand1012/parallel-stock-processor/files"
	"github.com/chand1012/parallel-stock-processor/stocks"
)

var VALID_MODES = []string{"stddev", "avg", "max", "min", "gain"}

func main() {

	tickersFlag := flag.String("tickers", "", "The ticker symbols to process")
	mode := flag.String("mode", "", "The mode to process the data in")

	flag.Parse()

	// get NUM_THREADS from env
	num_threads := os.Getenv("NUM_THREADS")
	if num_threads == "" {
		num_threads = "1"
	}

	run_verbose := os.Getenv("VERBOSE")
	if run_verbose == "" {
		run_verbose = "false"
	}

	var tickers []string
	if *tickersFlag != "" {
		// split the tickers by comma
		tickers = strings.Split(*tickersFlag, ",")
	} else {
		tickers = []string{}
	}

	threads, err := strconv.ParseInt(num_threads, 10, 64)
	if err != nil {
		panic(err)
	}

	var verbose bool
	if run_verbose != "" {
		verbose = true
	} else {
		verbose = false
	}

	if verbose {
		fmt.Println("Executing with " + num_threads + " threads")
	}

	start := time.Now()
	csvData, err := files.GetAllFiles("data", int(threads))
	if err != nil {
		panic(err)
	}
	end := time.Now()

	if verbose {
		fmt.Println("Time to load all files: " + end.Sub(start).String())
	}

	start = time.Now()

	rowChan := make(chan []stocks.StockRow, len(tickers))
	var wg sync.WaitGroup

	for _, ticker := range tickers {
		wg.Add(1)
		go func(data [][]string, rowChan chan []stocks.StockRow) {
			defer wg.Done()
			row, err := stocks.ParseCSV(data)
			if err != nil {
				panic(err)
			}
			rowChan <- row
		}(csvData[ticker], rowChan)
	}

	wg.Wait()

	// get the data from the channel
	stockData := make(map[string][]stocks.StockRow)
	for i := 0; i < len(tickers); i++ {
		data := <-rowChan
		stockData[data[0].Ticker] = data
	}
	end = time.Now()

	if verbose {
		fmt.Println("Time to parse all files: " + end.Sub(start).String())
	}

	switch *mode {
	case "stddev":
		fmt.Println("Running in Standard Deviation mode")
		for _, ticker := range tickers {
			closings := stocks.GetAllClosings(ticker, stockData)
			// fmt.Println(closings)
			stddev, err := stats.StandardDeviation(closings)
			if err != nil {
				panic(err)
			}
			fmt.Println(ticker + " " + strconv.FormatFloat(stddev, 'f', 2, 64))
		}
	case "avg":
		fmt.Println("Running in Average mode")
		for _, ticker := range tickers {
			closings := stocks.GetAllClosings(ticker, stockData)
			// fmt.Println(closings)
			avg, err := stats.Mean(closings)
			if err != nil {
				panic(err)
			}
			fmt.Println(ticker + " " + strconv.FormatFloat(avg, 'f', 2, 64))
		}
	case "max":
		fmt.Println("Running in Max mode")
		for _, ticker := range tickers {
			closings := stocks.GetAllClosings(ticker, stockData)
			// fmt.Println(closings)
			max, err := stats.Max(closings)
			if err != nil {
				panic(err)
			}
			fmt.Println(ticker + " " + strconv.FormatFloat(max, 'f', 2, 64))
		}
	case "min":
		fmt.Println("Running in Min mode")
		for _, ticker := range tickers {
			closings := stocks.GetAllClosings(ticker, stockData)
			// fmt.Println(closings)
			min, err := stats.Min(closings)
			if err != nil {
				panic(err)
			}
			fmt.Println(ticker + " " + strconv.FormatFloat(min, 'f', 2, 64))
		}
	case "gain":
		// gain = (closing - opening) / opening
		// then get the average per day gain
		fmt.Println("Running in Gain mode")
		for _, ticker := range tickers {
			openings := stocks.GetAllOpenings(ticker, stockData)
			closings := stocks.GetAllClosings(ticker, stockData)
			var gains []float64
			for i := 0; i < len(openings); i++ {
				gain := (closings[i] - openings[i]) / openings[i]
				gains = append(gains, gain)
			}
			avg, err := stats.Mean(gains)
			if err != nil {
				panic(err)
			}
			fmt.Println(ticker + " " + strconv.FormatFloat(avg, 'f', 64, 64))
		}
	}
}
