package files

import (
	"fmt"
	"io/ioutil"
	"sync"
)

func loadFile(fileName string, tickerName chan string, fileData chan [][]string, errs chan error) {
	// defer wg.Done()
	// load the file
	data, err := LoadCSV(fileName)
	if fileName == "" {
		return
	}

	if err != nil {
		fmt.Println("Error loading file: " + fileName)
		errs <- err
		return
	}

	// remove .csv from the file name
	ticker := fileName[5 : len(fileName)-4]

	// send the data to the channel
	fileData <- data
	tickerName <- ticker
}

func loadFileWorker(fileNames []string, tickerName chan string, fileData chan [][]string, errs chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, fileName := range fileNames {
		loadFile(fileName, tickerName, fileData, errs)
	}
}

func GetAllFiles(dir string, threads int) (map[string][][]string, error) {
	var wg sync.WaitGroup

	// get all the files in the directory
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	// create a channel to send the file data to
	fileData := make(chan [][]string, len(files))
	tickerName := make(chan string, len(files))

	// create a channel to send the errors to
	errs := make(chan error, len(files))

	fileNameChunks := [][]string{}
	for i := 0; i < threads; i++ {
		strLen := len(files) / threads
		s := make([]string, strLen)
		fileNameChunks = append(fileNameChunks, s)
	}

	currentThread := 0

	for _, file := range files {
		if currentThread >= threads {
			currentThread = 0
		}
		fileNameChunks[currentThread] = append(fileNameChunks[currentThread], dir+"/"+file.Name())
		currentThread++
	}

	for _, fileNames := range fileNameChunks {
		wg.Add(1)
		go loadFileWorker(fileNames, tickerName, fileData, errs, &wg)
	}

	// wait for all the files to be loaded
	wg.Wait()

	// close the channels
	close(fileData)
	close(errs)
	close(tickerName)

	// check if there were any errors
	if len(errs) > 0 {
		return nil, <-errs
	}

	// create a map to hold the data
	dataMap := make(map[string][][]string)

	// loop through the data and add it to the map
	for i := 0; i < len(files); i++ {
		dataMap[<-tickerName] = <-fileData
	}

	return dataMap, nil
}
