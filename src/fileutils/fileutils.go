// Package fileutils any file based interactions
package fileutils

import (
	"fmt"
	"os"
	"strings"

	"github.com/turnage/graw/reddit"
)

// WriteTextOutput to write txt output based on a channel to a file
// fileName: name of write to write to, c: channel
func WriteTextOutput(fileName string, c chan reddit.Comment) {
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

// HeadersToRegex simple function to convert a slice of string headers to a regex
/* headers: slice of headers to use in a "form style" regex, for example:
Headers: HEADER1,HEADER2,HEADER3
Regex returned: `HEADER1: (.*)\n\nHEADER2: (.*)\n\HEADER3: (.*)`
*/
func HeadersToRegex(headers []string) string {
	var b strings.Builder
	for _, header := range headers {
		fmt.Fprintf(&b, "%s: (.*)\\n\\n", header)
	}
	regString := b.String()
	return regString[:len(regString)-4] // Remove last newline char for form ending
}
