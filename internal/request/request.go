package request

import (
	"fmt"
	"io"
	"strings"
)

const bufferSize = 8

type parserState int

const (
	initialized = iota + 1
	done
)

type Request struct {
	RequestLine RequestLine
	state       parserState //0 = initialized, 1 = done
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	// read from reader in chunks and try to parse the request line until we have enough data to parse it, then return the request

	buf := make([]byte, bufferSize)

	readToIndex := 0 // how much data has been read from io.buffer

	currentReq := &Request{
		state: initialized,
	}
	for {
		// check if buffer is full
		if readToIndex >= len(buf) {
			// buffer is full, we need to grow it
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		n, err := reader.Read(buf[readToIndex:])
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n == 0 {
			break
		}
		readToIndex += n
		consumedBytes, err := currentReq.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}
		if consumedBytes > 0 {
			copy(buf, buf[consumedBytes:readToIndex])
			readToIndex -= consumedBytes
		}
		if currentReq.state == done {
			break
		}
	}

	return currentReq, nil
}

func parseRequestLine(line string) (*RequestLine, int, error) {
	//WRONG!: consumedBytes := len(line)

	//case needs more data to parse request line
	if !strings.Contains(line, "\r\n") {
		return nil, 0, nil
	}
	consumedBytes := 0
	for consumedBytes < len(line) {
		if line[consumedBytes] == '\r' {
			if consumedBytes+1 < len(line) && line[consumedBytes+1] == '\n' {
				break
			}
		}
		consumedBytes++
	}
	split := strings.Split(line, "\r\n")

	reqLine := split[0] // ignore for now

	partsOfReqLine := strings.Split(reqLine, " ")
	if len(partsOfReqLine) != 3 {
		return nil, 0, fmt.Errorf("invalid number of parts in request line: %d", len(partsOfReqLine))
	}
	// assign parts of req line
	method := partsOfReqLine[0]
	requestTarget := partsOfReqLine[1]
	httpVersion := partsOfReqLine[2]

	//check for valid data
	// method should be only capital letters
	for _, r := range method {
		if r < 'A' || r > 'Z' {
			return nil, 0, fmt.Errorf("invalid method: %q", method)
		}
	}
	if method == "" {
		return nil, 0, fmt.Errorf("method cannot be empty")
	}

	// httpversion should be HTTP/1.1
	if httpVersion != "HTTP/1.1" {
		return nil, 0, fmt.Errorf("invalid http version: %s", httpVersion)
	} else {
		httpVersion = "1.1"
	}

	return &RequestLine{
		HttpVersion:   httpVersion,
		RequestTarget: requestTarget,
		Method:        method,
	}, consumedBytes, nil

}

func (r *Request) parse(data []byte) (int, error) {
	if r.state == done {
		return 0, fmt.Errorf("trying to read data into a 'done' state: %v", data)
	} else if r.state == initialized {
		// parse request line
		// if we don't have enough data to parse the request line, return 0 and wait for more data
		// if we have enough data to parse the request line, parse it and update the state to done
		reqLine, consumedBytes, err := parseRequestLine(string(data))
		if err != nil {
			return 0, err
		}
		if reqLine == nil {
			return 0, nil
		}
		r.RequestLine = *reqLine
		r.state = done
		return consumedBytes, nil
	} else {
		return 0, fmt.Errorf("invalid state: %d", r.state)
	}
}
