package response

import (
	"fmt"
	"http_from_tcp/internal/headers"
	"io"
	"log"
)

type StatusCode int

const(
	StatusOK            StatusCode = 200
	StatusBadRequest    StatusCode = 400
	StatusInternalError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	var err error
	switch statusCode {
	case StatusOK:
		_, err = fmt.Fprintf(w, "HTTP/1.1 %d OK\r\n", statusCode)
	case StatusBadRequest:
		_, err = fmt.Fprintf(w, "HTTP/1.1 %d Bad Request\n", statusCode)
	case StatusInternalError:
		_, err = fmt.Fprintf(w, "HTTP/1.1 %d Internal Server Error\n", statusCode)
	default:
		_, err = fmt.Fprintf(w, "HTTP/1.1 %d", statusCode)
	}
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers["Content-Length"] = fmt.Sprintf("%d", contentLen)
	headers["Connection"] = "close"
	headers["Content-Type"] = "plain/text"
	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	var err error
	for k, v := range headers {
		_, err = fmt.Fprintf(w, "%s: %s\n", k, v)
	}
	_, err = fmt.Fprint(w, "\n")
	if err != nil {
		log.Fatalf("Error writing headers: %s\n", err)
	}
	return err
}

