package request

import (
	"HTTP_OVER_TCP/helpers"
	"errors"
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

func parseRequestLine(line string) (*RequestLine, error) {
	reqLine  := RequestLine{}
	segments := strings.Split(line, " ")
	if len(segments) < 3 {
		return &reqLine, errors.New("not enough parts in request line")
	}
	method, err := checkRequestMethod(segments[0])
	if err != nil {
		return nil, err
	}
	target, err := checkRequestTarget(segments[1])
	if err != nil {
		return nil, err
	}
	version, err := checkHttpVersion(segments[2])
	if err != nil {
		return nil, err
	}
	

	reqLine.Method        = method
	reqLine.RequestTarget = target
	reqLine.HttpVersion   = version
	return &reqLine, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	content, err := io.ReadAll(reader)
	helpers.CheckFatal(err, "Error reading content.")
	req := Request{}
	lines := strings.Split(string(content), "\r\n")
	reqLine, err := parseRequestLine(lines[0])
	if err != nil {
		return nil, err
	}
	// for line := range strings.SplitSeq(string(content), "\r\n") {
	// 	reqLine := parseRequestLine(line)
	req.RequestLine = *reqLine
	//
	// }
	return &req, nil
}
