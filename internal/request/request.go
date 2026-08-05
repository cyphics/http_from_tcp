package request

import (
	"errors"
	"fmt"
	"http_from_tcp/internal/headers"
	"io"
	"log"
	"strconv"
	"strings"
)

const bufferSize int = 8

type parserState string
const(
	StateInit    parserState = "StateInit"
	StateHeaders parserState = "StateHeaders"
	StateBody    parserState = "StateBody"
	StateDone    parserState = "StateDone"
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
	state       parserState
}

type RequestLine struct {
	HTTPVersion   string
	RequestTarget string
	Method        string
}

func buildRequest() *Request {
	req := Request{}
	req.Headers = headers.NewHeaders()
	req.state = StateInit
	return &req
}

func checkRequestMethod(method string) (string, error) {
	validMethods := [...]string{"GET", "POST"}
	for _, m := range validMethods {
		if m == method {
			return method, nil
		}
	}
	return method, errors.New("unknown request method")
}

func checkRequestTarget(target string) (string, error){
	if target[0] == '/' {
		return target, nil
	}
	return target, errors.New("request target should should start with a '/'")
}

func checkHttpVersion(version string) (string, error) {
	segs := strings.Split(version, "HTTP/")
	if len(segs) < 2 {
		return version, errors.New("bad HTTP version format")
	}
	return segs[1], nil
}

func (r *Request) parseRequestLine(data []byte) (int, error) {
	line, _, found := strings.Cut(string(data), "\r\n")
	if !found {
		return 0, nil
	}
	// fmt.Printf("Parsing request line '%s'\n", line)
	reqLine  := RequestLine{} 
	segments := strings.Split(line, " ")
	if len(segments) < 3 {
		return 0, fmt.Errorf("not enough parts in request line \"%s\"", line)
	}
	method, err := checkRequestMethod(segments[0])
	if err != nil { return 0, err }
	target, err := checkRequestTarget(segments[1])
	if err != nil { return 0, err }
	version, err := checkHttpVersion(segments[2])
	if err != nil { return 0, err }
	
	reqLine.Method        = method
	reqLine.HTTPVersion   = version
	reqLine.RequestTarget = target
	r.RequestLine   = reqLine
	r.state = StateHeaders
	return len(line)+2, nil
}

func (r *Request) parseBody(data []byte) (int, error) {
	r.Body = append(r.Body, data...)
	return len(data), nil
}

func (r *Request) parse(data []byte) (int, error) {
	// fmt.Printf("parsing chunk '%s'\n", data)
	var err error
	totalBytesRead := 0
	bytesRead      := 0
	done           := false

	for r.state != StateDone {

		switch r.state {
		case StateInit:
			bytesRead, err = r.parseRequestLine(data[totalBytesRead:])
			totalBytesRead += bytesRead
		case StateHeaders:
			bytesRead, done, err = r.Headers.Parse(data[totalBytesRead:])
			totalBytesRead += bytesRead
			if done {
				r.state = StateBody
				_, hasBody := r.Headers.Get("content-length")
				if hasBody {
					r.state = StateBody
				} else {
					r.state = StateDone
				}
			}
		case StateBody:
			bytesRead, err = r.parseBody(data[totalBytesRead:])
			totalBytesRead += bytesRead
			return totalBytesRead, err
		default:
			return 0, fmt.Errorf("parsing state error")
		}

		if err != nil { return totalBytesRead, err }
		if bytesRead == 0 {
			return totalBytesRead, nil
		}
	}
	return totalBytesRead, nil
}


func RequestFromReader(reader io.Reader) (*Request, error) {
	loopCounter := 0
	req         := buildRequest()
	readBuffer  := make([]byte, bufferSize)
	eof         := false
	var readIndex = 0

	for req.state != StateDone {

		loopCounter += 1
		if loopCounter > 100 {
			log.Fatalf("Infinite loop!")
		}

		// Read into the buffer
		numBytesRead, err := reader.Read(readBuffer[readIndex:])
		readIndex += numBytesRead
		if readIndex >= len(readBuffer) {
			tmpBuffer := make([]byte, 2 * len(readBuffer))
			copy(tmpBuffer, readBuffer)
			readBuffer = tmpBuffer
		}
		if err != nil {
			// If error is EOF, parse before treating the result
			if err == io.EOF {
				eof = true
			} else {
				return nil, err
			}
		}

		// Parse from the buffer
		numBytesParsed, err := req.parse(readBuffer[:readIndex])
		if err != nil { return nil, err }
		if numBytesParsed > 0 {
			readIndex -= numBytesParsed
			tmpBuffer := make([]byte, len(readBuffer))
			copy(tmpBuffer, readBuffer[numBytesParsed:])
			readBuffer = tmpBuffer
		}
		if eof {
			if req.state == StateDone { return req, nil }
			if req.state == StateHeaders { return req, errors.New("incomplete headers") }
			if req.state == StateBody {
				contentLengthStr,_ := req.Headers.Get("content-length")
				contentLength, err := strconv.Atoi(contentLengthStr)
				if err != nil { return req, err }
				if contentLength > len(req.Body) {
					return req, errors.New("incomplete body")
				}
				if contentLength < len(req.Body) {
					return req, errors.New("body longer than expected")
				}
				return req, nil
			}
			return req, errors.New("state error")
		}
	}
	return req, nil
}
