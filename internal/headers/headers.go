package headers

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	CRLF_Index := bytes.Index(data, []byte("\r\n"))

	if CRLF_Index == -1 {
		// No CRLF found, which means we don't have a complete header line yet
		return 0, false, nil
	}
	if CRLF_Index == 0 {
		// Found CRLF at the start, which indicates the end of headers
		return 2, true, nil
	}

	// We have a complete header line, so we can parse it
	line := string(data[:CRLF_Index])
	fmt.Println("Incoming data: ", line)

	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return 0, false, errors.New("invalid header format: " + line)
	}

	keyUntrimmed := parts[0]
	fmt.Printf("Parsing header line: '%s'\n", line)
	fmt.Printf(" - Key before trim: '%s'\n", keyUntrimmed)

	keyTrimmed := strings.TrimSpace(keyUntrimmed)
	fmt.Printf(" - Key after trim: '%s'\n", keyTrimmed)

	if len(keyTrimmed) != len(keyUntrimmed) {
		fmt.Printf(" - Key had leading/trailing spaces: '%s'\n", keyUntrimmed)
		return 0, false, errors.New("invalid header format: key has leading/trailing spaces: " + keyUntrimmed)
	}

	if !fieldNameCheck(keyTrimmed) {
		fmt.Printf(" - Key contains invalid characters: '%s'\n", keyTrimmed)
		return 0, false, errors.New("invalid header format: key contains invalid characters: " + keyTrimmed)
	}

	keyTrimmed = strings.ToLower(keyTrimmed) // Normalize header keys to lowercase
	value := strings.TrimSpace(parts[1])

	// if existing header, append to value with comma separation
	if preexistingKey, exists := h[keyTrimmed]; exists {
		value = preexistingKey + ", " + value
	}

	h[keyTrimmed] = value

	n += len(line) + 2 // +2 for the CRLF
	fmt.Printf("Parsed header: '%s: %s bytes: %d'\n", keyTrimmed, value, n)

	fmt.Printf("Finished parsing header. Total bytes parsed: %d\n", n)

	return n, false, err
}

func fieldNameCheck(fieldName string) bool {
	if len(fieldName) == 0 {
		fmt.Printf(" - Key is empty\n")
		return false
	}

	const validFieldNameChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!,#$%&'*+-.^_`|~"

	for _, r := range fieldName {
		if !strings.ContainsRune(validFieldNameChars, r) {
			fmt.Printf(" - Key contains invalid character: '%c'\n", r)
			return false
		}
	}
	return true
}
