package main

import (
	"fmt"

	"github.com/chand1012/parallel-term-project/files"
	"github.com/chand1012/parallel-term-project/util"
)

func main() {
	csvData, err := files.LoadCSV("data/GOOG.csv")
	if err != nil {
		panic(err)
	}

	data, err := util.ParseCSV(csvData)
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
