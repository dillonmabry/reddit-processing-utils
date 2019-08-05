// Package fileutils any file based interactions
package fileutils

import (
	"fmt"
	"os"

	graw "github.com/turnage/graw/reddit"
)

// WriteTextOutput to write txt output based on a channel to a file
// fileName: name of write to write to, c: channel
func WriteTextOutput(fileName string, c chan graw.Comment) {
	fname := fmt.Sprintf("%s%s", fileName, ".txt")
	f, err := os.OpenFile(fname, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for msg := range c {
		if _, err := f.WriteString(msg.Body + "\n"); err != nil {
			panic(err)
		}
	}
}

// WriteCsvOutput to write csv output based on a channel to a file
// fileName: name of write to write to, c: channel
func WriteCsvOutput(fileName string, headers []string, c chan []string) {
	fname := fmt.Sprintf("%s%s", fileName, ".csv")
	w, err := NewCsvWriter(fname)
	if err != nil {
		panic(err)
	}
	w.Write(headers)

	for msg := range c {
		w.Write(msg)
	}
	defer w.Flush()
}
