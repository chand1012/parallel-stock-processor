package main

import (
	"fmt"

	"github.com/chand1012/parallel-stock-processor/files"
	"github.com/chand1012/parallel-stock-processor/stocks"
)

func main() {
	csvData, err := files.GetAllFiles("data")
	if err != nil {
		panic(err)
	}

	for _, file := range csvData {
		data, err := stocks.ParseCSV(file)
		if err != nil {
			panic(err)
		}

		totalClosingPrice := float64(0)

		for _, row := range data {
			totalClosingPrice += row.Close
		}

		// print the average closing price
		fmt.Println(totalClosingPrice / float64(len(data)))
	}
}
