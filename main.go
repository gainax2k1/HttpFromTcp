package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	msgs, err := os.Open("messages.txt")
	if err != nil {
		panic(err)
	}
	defer msgs.Close()

	msgBuffer := make([]byte, 8)
	var currentLine string

	for bytesRead, err := msgs.Read(msgBuffer); err != io.EOF; bytesRead, err = msgs.Read(msgBuffer) {
		if err != nil {
			panic(err)
		}

		currentLine += string(msgBuffer[:bytesRead])

		//check for newlines in currentLine
		if strings.Contains(currentLine, "\n") {
			//split currentLine by newlines
			splitStrings := strings.Split(currentLine, "\n")

			for i, line := range splitStrings {
				//  print each line except the last one, which may be incomplete
				if i == len(splitStrings)-1 {
					// last line, may be incomplete, save it for the next read
					currentLine = line
				} else {
					fmt.Printf("read: %s\n", line)
				}
			}

		}
		// old output
		//fmt.Printf("read: %s\n", string(msgBuffer[:bytesRead]))
	}
}
