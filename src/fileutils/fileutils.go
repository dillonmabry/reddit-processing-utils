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
	f, err := os.OpenFile(fmt.Sprintf("%s%s", fileName, ".txt"),
		os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for msg := range c {
		f.WriteString(msg.Body + "\n")
	}
}
