package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func main() {
	/* old method from file
	msgs, err := os.Open("messages.txt")
	if err != nil {
		panic(err)
	}
	defer msgs.Close()
	*/
	const port = ":42069"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}
	defer listener.Close()

	fmt.Println("Listening on", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		fmt.Printf("Accepted connection from %s\n", conn.RemoteAddr())

		linesChan := getLinesChannel(conn)
		for line := range linesChan {
			fmt.Println(line)
		}

		fmt.Println("Closed connection from", conn.RemoteAddr())
	}

}

// lines is a channel that will receive lines read from the file

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
							// Send the line to the channel
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
			// Send any remaining text in currentLine to the channel
			lines <- currentLine
		}

	}()

	return lines
}
