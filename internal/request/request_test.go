package request

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

// Read reads up to len(p) or numBytesPerRead bytes from the string per call
// its useful for simulating reading a variable number of bytes per chunk from a network connection
func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := cr.pos + cr.numBytesPerRead
	if endIndex > len(cr.data) {
		endIndex = len(cr.data)
	}
	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n

	return n, nil
}

func TestRequestLineParse(t *testing.T) {
	// Test: Direct use of parseRequestLine()
	req := Request{}
	data :=  "GET / HTTP/1.1\r\nHost"
	b, e := req.parseRequestLine([]byte(data))
	require.NoError(t, e)
	require.NotNil(t, req)
	require.Equal(t, b, 16)
	assert.Equal(t, "GET", req.RequestLine.Method)
	assert.Equal(t, "/", req.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", req.RequestLine.HTTPVersion)

	// Test: Good GET Request line
	data = "GET / HTTP/1.1\r\n" +
		"Host: localhost:42069\r\n" +
		"User-Agent: curl/7.81.0\r\n" +
		"Accept: */*\r\n" +
		"\r\n"
	r, err := RequestFromReader(strings.NewReader(data))
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HTTPVersion)

	// Test: Good GET Request line with path
	data = "GET /coffee HTTP/1.1\r\n" +
		"Host: localhost:42069\r\n" +
		"User-Agent: curl/7.81.0\r\n" +
		"Accept: */*\r\n" +
		"\r\n"
	r, err = RequestFromReader(strings.NewReader(data))
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HTTPVersion)

	// Test: Good POST Request line with path
	r, err = RequestFromReader(strings.NewReader("POST /coffee HTTP/1.1\r\n" +
		"Host: localhost:42069\r\n" +
		"User-Agent: curl/7.81.0\r\n" +
		"Accept: */*\r\n" +
		"\r\n" +
		""))
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "POST", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HTTPVersion)

	// Test: Invalid number of parts in request line
	_, err = RequestFromReader(strings.NewReader("/coffee HTTP/1.1\r\n" +
		"Host: localhost:42069\r\n" +
		"User-Agent: curl/7.81.0\r\n" +
		"Accept: */*\r\n" +
		"\r\n"))
	require.Error(t, err)

	// Test: Invalid method (out of order) Request line
	_, err = RequestFromReader(strings.NewReader("/coffee GET HTTP/1.1\r\n" +
		"Host: localhost:42069\r\n" +
		"User-Agent: curl/7.81.0\r\n" +
		"Accept: */*\r\n" +
		"\r\n"))
	require.Error(t, err)

	// Test: Invalid method (out of order) Request line part 2
	_, err = RequestFromReader(strings.NewReader("GET HTTP/1.1 /coffee\r\n" +
		"Host: localhost:42069\r\n" +
		"User-Agent: curl/7.81.0\r\n" +
		"Accept: */*\r\n" +
		"\r\n"))
	require.Error(t, err)

	// Test: Invalid HTTP version
	_, err = RequestFromReader(strings.NewReader("GET /coffee 1.1\r\n" +
		"Host: localhost:42069\r\n" +
		"User-Agent: curl/7.81.0\r\n" +
		"Accept: */*\r\n" +
		"\r\n"))
	require.Error(t, err)
}

func TestChunkReader(t *testing.T) {
	// Test: Good GET Request line
	data := "GET / HTTP/1.1\r\n" +
		"Host: localhost:42069\r\n" +
		"User-Agent: curl/7.81.0\r\n" +
		"Accept: */*\r\n" +
		"\r\n"
	reader := &chunkReader{
		data:            data,
		numBytesPerRead: 3,
	}
	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HTTPVersion)

	// Test: Good GET Request line with path
	data = "GET /coffee HTTP/1.1\r\n" +
		"Host: localhost:42069\r\n" +
		"User-Agent: curl/7.81.0\r\n" +
		"Accept: */*\r\n" +
		"\r\n"
	for i := range len(data) - 1 {
		reader := &chunkReader{
			data:            data,
			numBytesPerRead: i + 1,
		}
		r, err := RequestFromReader(reader)
		require.NoError(t, err)
		require.NotNil(t, r)
		assert.Equal(t, "GET", r.RequestLine.Method)
		assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
		assert.Equal(t, "1.1", r.RequestLine.HTTPVersion)
	}

	// Test: Invalid number of parts in request line
	_, err = RequestFromReader(&chunkReader{
		data: "/coffee HTTP/1.1\r\n" +
			"Host: localhost:42069\r\n" +
			"User-Agent: curl/7.81.0\r\n" +
			"Accept: */*\r\n" +
			"\r\n",
		numBytesPerRead: 3,
	})
	require.Error(t, err)

	// Test: Invalid method (out of order) Request line
	_, err = RequestFromReader(&chunkReader{
		data: "/coffee GET HTTP/1.1\r\n" +
			"Host: localhost:42069\r\n" +
			"User-Agent: curl/7.81.0\r\n" +
			"Accept: */*\r\n" +
			"\r\n",
		numBytesPerRead: 3,
	})
	require.Error(t, err)

	// Test: Invalid method (out of order) Request line part 2
	_, err = RequestFromReader(&chunkReader{
		data: "GET HTTP/1.1 /coffee\r\n" +
			"Host: localhost:42069\r\n" +
			"User-Agent: curl/7.81.0\r\n" +
			"Accept: */*\r\n" +
			"\r\n",
		numBytesPerRead: 3,
	})
	require.Error(t, err)

	// Test: Invalid HTTP version
	_, err = RequestFromReader(&chunkReader{
		data: "GET /coffee 1.1\r\n" +
			"Host: localhost:42069\r\n" +
			"User-Agent: curl/7.81.0\r\n" +
			"Accept: */*\r\n" +
			"\r\n",
		numBytesPerRead: 3,
	})
	require.Error(t, err)
}

func TestRequestHeader(t *testing.T) {
	// Test: Standard Headers
	reader := &chunkReader{
		data: "GET / HTTP/1.1\r\n" +
			"Host: localhost:42069\r\n" +
			"User-Agent: curl/7.81.0\r\n" +
			"Accept: */*\r\n" +
			"\r\n",
		numBytesPerRead: 9,
	}
	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "localhost:42069", r.Headers["host"])
	assert.Equal(t, "curl/7.81.0", r.Headers["user-agent"])
	assert.Equal(t, "*/*", r.Headers["accept"])

	// Test: Malformed Header
	reader = &chunkReader{
		data: "GET / HTTP/1.1\r\n" +
			"Host localhost:42069\r\n" +
			"\r\n",
		numBytesPerRead: 3,
	}
	_, err = RequestFromReader(reader)
	require.Error(t, err)

	// Test: Empty headers
	reader = &chunkReader{
		data:            "GET / HTTP/1.1\r\n\r\n",
		numBytesPerRead: 3,
	}
	r, err = RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "", r.Headers["host"])

	// Test: Duplicate headers
	reader = &chunkReader{
		data: "GET / HTTP/1.1\r\n" +
			"Set-Person: lane-loves-go\r\n" +
			"Set-Person: prime-loves-zig\r\n" +
			"Set-Person: tj-loves-ocaml \r\n" +
			"Host: localhost:42069\r\n" +
			"User-Agent: curl/7.81.0\r\n" +
			"Accept: */*\r\n" +
			"\r\n",
		numBytesPerRead: 3,
	}
	r, err = RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "localhost:42069", r.Headers["host"])
	assert.Equal(t, "curl/7.81.0", r.Headers["user-agent"])
	assert.Equal(t, "*/*", r.Headers["accept"])
	assert.Equal(t, "*/*", r.Headers["accept"])
	assert.Equal(t, "lane-loves-go; prime-loves-zig; tj-loves-ocaml", r.Headers["set-person"])

	// Test: case insensitive headers
	reader = &chunkReader{
		data: "GET / HTTP/1.1\r\n" +
			"set-person: lane-loves-go\r\n" +
			"Set-Person: prime-loves-zig\r\n" +
			"Set-Person: tj-loves-ocaml \r\n" +
			"Host: localhost:42069\r\n" +
			"User-Agent: curl/7.81.0\r\n" +
			"Accept: */*\r\n" +
			"\r\n",
		numBytesPerRead: 3,
	}
	r, err = RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "localhost:42069", r.Headers["host"])
	assert.Equal(t, "curl/7.81.0", r.Headers["user-agent"])
	assert.Equal(t, "*/*", r.Headers["accept"])
	assert.Equal(t, "*/*", r.Headers["accept"])
	assert.Equal(t, "lane-loves-go; prime-loves-zig; tj-loves-ocaml", r.Headers["set-person"])

	// Test: Missing end of header
	reader = &chunkReader{
		data: "GET / HTTP/1.1\r\n"+
		"Host: localhost:42069\r\n",
		numBytesPerRead: 3,
	}
	_, err = RequestFromReader(reader)
	require.Error(t, err)
}

func TestRequestBody(t *testing.T) {
	// Test: Standard Body
	reader := &chunkReader{
		data: "POST /submit HTTP/1.1\r\n" +
			"Host: localhost:42069\r\n" +
			"Content-Length: 13\r\n" +
			"\r\n" +
			"hello world!\n",
		numBytesPerRead: 3,
	}
	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "hello world!\n", string(r.Body))

	// Test: Body shorter than reported content length
	reader = &chunkReader{
		data: "POST /submit HTTP/1.1\r\n" +
			"Host: localhost:42069\r\n" +
			"Content-Length: 20\r\n" +
			"\r\n" +
			"partial content",
		numBytesPerRead: 3,
	}
	_, err = RequestFromReader(reader)
	require.Error(t, err)
}
