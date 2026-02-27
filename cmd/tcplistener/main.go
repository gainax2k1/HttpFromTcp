package main

import (
	"fmt"
	"log"
	"net"

	"HttpFromTcp/internal/request"
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

		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Printf("Error reading request: %s\n", err)
		} else {
			fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion)
		}

		fmt.Println("Closed connection from", conn.RemoteAddr())
	}

}
