package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	msgs, err := os.Open("messages.txt")
	if err != nil {
		panic(err)
	}
	defer msgs.Close()

	msgBuffer := make([]byte, 8)

	for bytesRead, err := msgs.Read(msgBuffer); err != io.EOF; bytesRead, err = msgs.Read(msgBuffer) {
		fmt.Printf("read: %s\n", string(msgBuffer[:bytesRead]))
	}
}
