package files

import (
	"io/ioutil"
	"sync"
)

func loadFile(fileName string, fileData chan [][]string, errs chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	// load the file
	data, err := LoadCSV(fileName)
	if err != nil {
		errs <- err
		return
	}

	// send the data to the channel
	fileData <- data
}

func GetAllFiles(dir string) ([][][]string, error) {
	var wg sync.WaitGroup

	// get all the files in the directory
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	// create a channel to send the file data to
	fileData := make(chan [][]string, len(files))

	// create a channel to send the errors to
	errs := make(chan error, len(files))

	for _, file := range files {
		wg.Add(1)
		go loadFile(dir+"/"+file.Name(), fileData, errs, &wg)
	}

	// wait for all the files to be loaded
	wg.Wait()

	// close the channels
	close(fileData)
	close(errs)

	// check if there were any errors
	if len(errs) > 0 {
		return nil, <-errs
	}

	// create a slice to hold all the file data
	var allData [][][]string

	// read the data from the channel
	for data := range fileData {
		allData = append(allData, data)
	}

	return allData, nil
}
