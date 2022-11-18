package files

import (
	"encoding/csv"
	"fmt"
	"os"
)

func LoadCSV(filePath string) ([][]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println(filePath)
		fmt.Println(err)
		return nil, err
	}

	return records, err
}
