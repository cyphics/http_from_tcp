package response

import (
	"bytes"
	"errors"
	"fmt"
	"http_from_tcp/internal/headers"
	"log"
)

type StatusCode int
const(
	StatusOK            StatusCode = 200
	StatusBadRequest    StatusCode = 400
	StatusInternalError StatusCode = 500
)

type WriterState int
const(
	WriterStateStatusLine WriterState = iota
	WriterStateHeaders
	WriterStateBody
)

type Writer struct {
	Buffer bytes.Buffer
	state WriterState 
}

func (w *Writer) Write(data []byte) (int, error) {
	b, err := w.Buffer.Write(data)
	return b, err
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != WriterStateStatusLine {
		return errors.New("status line already written")
	}
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
	w.state = WriterStateHeaders
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers["Content-Length"] = fmt.Sprintf("%d", contentLen)
	headers["Connection"] = "close"
	headers["Content-Type"] = "plain/text"
	return headers
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state == WriterStateStatusLine {
		return errors.New("status line expected")
	}
	if w.state == WriterStateBody {
		return errors.New("body expected")
	}
	var err error
	for k, v := range headers {
		_, err = fmt.Fprintf(w, "%s: %s\n", k, v)
	}
	_, err = fmt.Fprint(w, "\n")
	if err != nil {
		log.Fatalf("Error writing headers: %s\n", err)
	}
	w.state = WriterStateBody
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != WriterStateBody {
		return 0, errors.New("body expected")
	}
	return w.Write(p)
}


