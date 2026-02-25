package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	reqString, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	parsedReqLine, err := parseRequestLine(string(reqString))
	if err != nil {
		return nil, err
	}
	return &Request{
		RequestLine: *parsedReqLine,
	}, nil
}

func parseRequestLine(line string) (*RequestLine, error) {
	split := strings.Split(line, "\r\n")

	reqLine := split[0] // ignore for now

	partsOfReqLine := strings.Split(reqLine, " ")
	if len(partsOfReqLine) != 3 {
		return nil, fmt.Errorf("invalid number of parts in request line: %d", len(partsOfReqLine))
	}
	// assign parts of req line
	method := partsOfReqLine[0]
	requestTarget := partsOfReqLine[1]
	httpVersion := partsOfReqLine[2]

	//check for valid data
	// method should be only capital letters
	for _, r := range method {
		if r < 'A' || r > 'Z' {
			return nil, fmt.Errorf("invalid method: %q", method)
		}
	}
	if method == "" {
		return nil, fmt.Errorf("method cannot be empty")
	}

	// httpversion should be HTTP/1.1
	if httpVersion != "HTTP/1.1" {
		return nil, fmt.Errorf("invalid http version: %s", httpVersion)
	} else {
		httpVersion = "1.1"
	}

	return &RequestLine{
		HttpVersion:   httpVersion,
		RequestTarget: requestTarget,
		Method:        method,
	}, nil

}
