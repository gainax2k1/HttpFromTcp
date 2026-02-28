package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Valid single header
	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid single header with extra whitespace
	headers = NewHeaders()
	data = []byte("Host:    localhost:42069    \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 30, n)
	assert.False(t, done)

	// Test: Valid 2 headers with existing headers
	headers = NewHeaders()
	headers["existing"] = "value"
	data = []byte("Host: localhost:42069\r\nUser-Agent: TestAgent\r\n\r\n")
	// - First call - parses Host
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)
	// - Second call - parses User-Agent (starting from where we left off)
	n2, done2, err2 := headers.Parse(data[n:])
	require.NoError(t, err2)
	assert.Equal(t, "TestAgent", headers["user-agent"])
	assert.Equal(t, 23, n2) // "User-Agent: TestAgent\r\n" = 23 bytes
	assert.False(t, done2)
	// - Third call - finds the empty line (end of headers)
	n3, done3, err3 := headers.Parse(data[n+n2:])
	require.NoError(t, err3)
	assert.Equal(t, 2, n3) // Just the final \r\n
	assert.True(t, done3)
	// - Existing header should still be there
	assert.Equal(t, "value", headers["existing"])

	// Test: Valid done parsing headers
	headers = NewHeaders()
	data = []byte("\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 2, n)
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid header format (missing colon)
	headers = NewHeaders()
	data = []byte("InvalidHeaderWithoutColon\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid header format (empty key)
	headers = NewHeaders()
	data = []byte(": NoKey\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: invalid fieldname characters
	headers = NewHeaders()
	data = []byte("Invalid@Header: value\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: invalid fieldname characters empty space
	headers = NewHeaders()
	data = []byte("Invalid Header: value\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: adding multiple values to the same header key
	headers = NewHeaders()
	data = []byte("Set-Cookie: cookie1=value1\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "cookie1=value1", headers["set-cookie"])
	assert.Equal(t, 28, n) // "Set-Cookie: cookie1=value1\r\n" = 28 bytes
	assert.False(t, done)
	data = []byte("Set-Cookie: cookie2=value2\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "cookie1=value1, cookie2=value2", headers["set-cookie"])
	assert.Equal(t, 28, n) // "Set-Cookie: cookie2=value2\r\n" = 28 bytes
	assert.False(t, done)
	n, done, err = headers.Parse([]byte("\r\n"))
	require.NoError(t, err)
	assert.Equal(t, 2, n) // Just the final \r\n
	assert.True(t, done)

}
