package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chand1012/parallel-stock-processor/files"
	"github.com/chand1012/parallel-stock-processor/stocks"
)

var VALID_MODES = []string{"stddev", "avg", "max", "min", "sum", "gain"}

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
	stocks.ParseMap(csvData, int(threads))
	end = time.Now()

	if verbose {
		fmt.Println("Time to parse all files: " + end.Sub(start).String())
	}

	if len(tickers) == 0 {
		for t := range csvData {
			tickers = append(tickers, t)
		}
	}

	switch *mode {
	case "stddev":
		start = time.Now()

	}

}
