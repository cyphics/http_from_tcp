package headers

import (
	"errors"
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	parsedBytes := 0
	done = false
	fieldLines := strings.Split(string(data), "\r\n")
	if len(fieldLines) < 2 {
		return 0, false, nil
	}
	if len(fieldLines[0]) == 0 {
		fmt.Println("done")
		done = true
	}
	for _, fieldLine := range fieldLines {
		parsedBytes += len(fieldLine) // accounts for \r\n delimiter
		if len(fieldLine) > 0 {
			if fieldLine[0] == ' ' {
				return 0, done, errors.New("whitespace not allowed at beginning of line")
			}
			key, value, _ := strings.Cut(fieldLine, ":")
			if key[len(key)-1] == ' ' {
				return 0, done, errors.New("whitespace not allowed between key and ':'")
			}
			h[strings.ToLower(key)] = strings.Trim(value, " ")
		}
	}
	return parsedBytes + len(fieldLines) - 1, done, nil
}
