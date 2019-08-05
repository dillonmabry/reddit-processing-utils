package fileutils

import (
	"encoding/csv"
	"os"
	"sync"
)

// CsvWriter custom type for safe concurrent access for csv writing
// Source: https://markhneedham.com/blog/2017/01/31/go-multi-threaded-writing-csv-file/
type CsvWriter struct {
	mutex     *sync.Mutex
	csvWriter *csv.Writer
}

// NewCsvWriter creates new writer and returns reference
func NewCsvWriter(fileName string) (*CsvWriter, error) {
	csvFile, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	w := csv.NewWriter(csvFile)
	return &CsvWriter{csvWriter: w, mutex: &sync.Mutex{}}, nil
}

func (w *CsvWriter) Write(row []string) {
	w.mutex.Lock()
	w.csvWriter.Write(row)
	w.mutex.Unlock()
}

// Flush the file, lock, unlock mutex
func (w *CsvWriter) Flush() {
	w.mutex.Lock()
	w.csvWriter.Flush()
	w.mutex.Unlock()
}
