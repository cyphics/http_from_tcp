package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaders(t *testing.T) {
	// TEST: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// TEST: Valid single header with extra whitespaces
	headers = NewHeaders()
	data = []byte("Host:    localhost:42069  \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 28, n)
	assert.False(t, done)

	// TEST: Valid Second header from existing header
	data = []byte("User-Agent: curl v8.10  \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, "curl v8.10", headers["user-agent"])
	assert.Equal(t, 26, n)
	assert.False(t, done)

	// TEST: Valid Done
	data = []byte("\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, "curl v8.10", headers["user-agent"])
	assert.Equal(t, 2, n)
	assert.True(t, done)


	// TEST: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// TEST: Invalid spacing after head
	headers = NewHeaders()
	data = []byte("Host : localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)


	// TEST: Invalid character
	headers = NewHeaders()
	data = []byte("H@st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// TEST Valid multiple values for same key
	headers = NewHeaders()
	n, done, err = headers.Parse([]byte("Set-Person: lane-loves-go\r\n"))
	n, done, err = headers.Parse([]byte("Set-Person: prime-loves-zig\r\n"))
	n, done, err = headers.Parse([]byte("Set-Person: tj-loves-ocaml\r\n"))
	require.NoError(t, err)
	assert.Equal(t, 28, n)
	assert.False(t, done)
	assert.Equal(t, "lane-loves-go; prime-loves-zig; tj-loves-ocaml", headers["set-person"])
}
