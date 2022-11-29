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
	// This gets parsed to an int64 to determine
	// the number of threads to use
	num_threads := os.Getenv("NUM_THREADS")
	if num_threads == "" {
		num_threads = "1"
	}

	// Sets if we should be verbose or not
	run_verbose := os.Getenv("VERBOSE")
	if run_verbose == "" {
		run_verbose = "false"
	}

	// parses the comma separated tickers into an array
	var tickers []string
	if *tickersFlag != "" {
		// split the tickers by comma
		tickers = strings.Split(*tickersFlag, ",")
	} else {
		tickers = []string{}
	}

	// parses the number of threads into an int64
	threads, err := strconv.ParseInt(num_threads, 10, 64)
	if err != nil {
		panic(err)
	}

	// parses the verbose flag into a bool
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
	// loads all the files using the number of threads given.
	// Loads all files in the "data" directory
	csvData, err := files.GetAllFiles("data", int(threads))
	if err != nil {
		panic(err)
	}
	end := time.Now()

	if verbose {
		fmt.Println("Time to load all files: " + end.Sub(start).String())
	}

	start = time.Now()

	// creates a channel to pass the data to.
	// Basically a queue that works across threads
	rowChan := make(chan []stocks.StockRow, len(tickers))
	var wg sync.WaitGroup

	for _, ticker := range tickers {
		wg.Add(1) // adds a thread to the wait group
		// creates a thread to process the data
		go func(data [][]string, rowChan chan []stocks.StockRow) {
			// defer finishing the thread until after the function is done
			defer wg.Done()
			// processes the data
			row, err := stocks.ParseCSV(data)
			if err != nil {
				// crashes the program and all threads if there is an error
				panic(err)
			}
			// puts the data into the channel
			rowChan <- row
		}(csvData[ticker], rowChan) // executes the thread
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
