package reader

import (
	"encoding/csv"
	"os"
)

type CSVReader struct {
	*Reader
}

func init() {
	Register("csv", newCSVReader)
}

func newCSVReader(r *Reader) IReader {
	return &CSVReader{r}
}

func (r *CSVReader) Read() ([][]string, error) {
	file, err := os.Open(r.Path)
	if err != nil {
		return nil, errorTableReadFailed(r.Path, err)
	}

	defer func(file *os.File) { _ = file.Close() }(file)

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, errorTableReadFailed(r.Path, err)
	}

	if records == nil || len(records) == 0 {
		return nil, errorTableNotSheet(r.Path)
	}

	return records, nil
}
