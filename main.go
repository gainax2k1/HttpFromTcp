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

	// lines is a channel that will receive lines read from the file
	lines := getLinesChannel(msgs)

	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}

	/* original code without channel:
	msgBuffer := make([]byte, 8)
	var currentLine string

	for {
		bytesRead, err := msgs.Read(msgBuffer)

		if bytesRead > 0 {
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

		}

		if err != nil {
			if err == io.EOF {
				break // end of file reached, exit loop.
			}
			panic(err)
		}

	}
	// print any remaining text in currentLine after the loop (in case last line does not end with a newline)
	if len(currentLine) > 0 {
		fmt.Printf("read: %s\n", currentLine)
	}
	*/
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)

	go func() {
		// Close the channel when the goroutine exits
		defer close(lines)
		defer f.Close()

		buf := make([]byte, 8)
		var currentLine string

		for {
			bytesRead, err := f.Read(buf)
			if bytesRead > 0 {
				currentLine += string(buf[:bytesRead])
				if strings.Contains(currentLine, "\n") {
					splitStrings := strings.Split(currentLine, "\n")
					for i, line := range splitStrings {
						if i == len(splitStrings)-1 {
							currentLine = line
						} else {
							lines <- line
						}
					}

				}
			}

			if err != nil {
				if err == io.EOF {
					break // end of file reached, exit loop.
				}
				panic(err)
			}
		}
		if len(currentLine) > 0 {
			lines <- currentLine
		}

	}()

	return lines
}
