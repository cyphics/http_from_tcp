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
		_, err = fmt.Fprintf(w, "HTTP/1.1 %d Bad Request\r\n", statusCode)
	case StatusInternalError:
		_, err = fmt.Fprintf(w, "HTTP/1.1 %d Internal Server Error\r\n", statusCode)
	default:
		_, err = fmt.Fprintf(w, "HTTP/1.1 %d\r\n", statusCode)
	}
	w.state = WriterStateHeaders
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	if contentLen > 0 {
		headers["Content-Length"] = fmt.Sprintf("%d", contentLen)
	} else {
		delete(headers, "Content-Length")
	}
	headers["Connection"] = "close"
	headers["content-type"] = "text/plain"
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
		_, err = fmt.Fprintf(w, "%s: %s\r\n", k, v)
	}
	_, err = fmt.Fprint(w, "\r\n")
	if err != nil { log.Fatalf("Error writing headers: %s\n", err) }
	w.state = WriterStateBody
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != WriterStateBody {
		return 0, errors.New("body expected")
	}
	return w.Write(p)
}

func (w *Writer) WriteChunkedBody(p []byte) (int,error) {
	return fmt.Fprintf(w, "%x\r\n%s\r\n", len(p), p)
}

func (w *Writer) WriteChunkedBodyDone() (int,error) {
	fmt.Println("Body done")
	return fmt.Fprintf(w, "%x\r\n\r\n", 0)
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	
	return nil
}
