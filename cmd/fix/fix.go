package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	fileutils "github.com/chand1012/parallel-stock-processor/files"
)

func main() {
	// get all the files in the directory
	files, err := ioutil.ReadDir("data")
	if err != nil {
		panic(err)
	}

	var brokenFiles []string

	// loop through each file
	// and try to load it as a csv
	for _, file := range files {
		// load the file
		_, err := fileutils.LoadCSV("data/" + file.Name())
		if err != nil {
			brokenFiles = append(brokenFiles, file.Name())
		}
	}

	fmt.Printf("Found %d broken files\n", len(brokenFiles))

	// read each line, count the number of commas
	// if there are less than 5 commas, delete the line
	for _, file := range brokenFiles {
		// load the file
		data, err := os.Open("data/" + file)
		if err != nil {
			panic(err)
		}

		// split the file into lines
		var lines []string
		scanner := bufio.NewScanner(data)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			panic(err)
		}

		var newLines []string

		for _, line := range lines {
			// count the number of commas
			commas := strings.Count(line, ",")
			if commas == 5 {
				newLines = append(newLines, line)
			}
		}

		f, err := os.Create("data/" + file + ".fixed")
		if err != nil {
			panic(err)
		}

		for _, line := range newLines {
			f.WriteString(line + "\n")
		}

		data.Close()
		f.Close()

		// delete the old file
		// and rename the new file

		err = os.Remove("data/" + file)
		if err != nil {
			panic(err)
		}

		err = os.Rename("data/"+file+".fixed", "data/"+file)
		if err != nil {
			panic(err)
		}
	}
}
