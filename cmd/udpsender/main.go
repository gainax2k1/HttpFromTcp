package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	const port = ":42069"

	UDPAddress, err := net.ResolveUDPAddr("udp", "localhost"+port)
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}

	conn, err := net.DialUDP("udp", nil, UDPAddress)
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}
	defer conn.Close()

	fmt.Println("Connected to", port)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("error: %s\n", err)
		}

		_, err = conn.Write([]byte(text))
		if err != nil {
			log.Fatalf("error: %s\n", err)
		}

	}
}
