package request

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
)

const bufferSize int = 8

type parserState int
const(
	Initialized parserState = iota
	Done
)

type Request struct {
	parsedData string
	RequestLine RequestLine
	state parserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func buildRequest() *Request {
	req := Request{}
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

func parseRequestLine(line string, request *Request) (int, error) {
	reqLine  := RequestLine{} 
	segments := strings.Split(line, " ")
	if len(segments) < 3 {
		return 0, fmt.Errorf("not enough parts in request line \"%s\"", line)
	}
	method, err := checkRequestMethod(segments[0])
	if err != nil {
		return 0, err
	}
	target, err := checkRequestTarget(segments[1])
	if err != nil {
		return 0, err
	}
	version, err := checkHttpVersion(segments[2])
	if err != nil {
		return 0, err
	}
	
	reqLine.Method        = method
	reqLine.HttpVersion   = version
	reqLine.RequestTarget = target
	request.RequestLine   = reqLine
	request.state = Done
	return len(line), nil
}

func (r *Request) parse(data []byte) (int, error) {
	lines := strings.Split(string(data), "\r\n")
	if len(lines) < 2 {
		return 0, nil
	}
	return parseRequestLine(lines[0], r)
}


func RequestFromReader(reader io.Reader) (*Request, error) {
	loopCounter := 0
	req := Request{}
	readBuffer := make([]byte, bufferSize, bufferSize)
	var readIndex = 0

	for req.state != Done {

		loopCounter += 1
		if loopCounter > 100 {
			log.Fatalf("Infinite loop!")
		}

		numBytesRead, err := reader.Read(readBuffer[readIndex:])
		readIndex += numBytesRead
		if readIndex >= len(readBuffer) {
			tmpBuffer := make([]byte, 2 * len(readBuffer), 2 * len(readBuffer))
			copy(tmpBuffer, readBuffer)
			readBuffer = tmpBuffer
		}
		if err != nil {
			if err == io.EOF {
				req.state = Done
				break
			}
			return nil, err
		}

		numBytesParsed, err := req.parse(readBuffer[:readIndex])
		if err != nil {
			return nil, err
		}
		if numBytesParsed > 0 {
			readIndex = 0
			readBuffer = readBuffer[numBytesParsed:]
		}
	}
	return &req, nil
}
