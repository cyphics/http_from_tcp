package headers

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func checkFieldName(name string) error {
	valid_chars := "!#$%&'*+-.^_`|~"
	for _, char := range name {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && !strings.Contains(valid_chars, string(char)) {
			return fmt.Errorf("character not allowed: %c in \"%s\"", char, name)
		}
	}
	return nil
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	fmt.Printf("Parses new header sequence \"%s\"\n\n", data)
	done = false
	fieldLines := strings.Split(string(data), "\r\n")
	if len(fieldLines) < 2 {
		return 0, false, nil
	}
	if len(fieldLines[0]) == 0 {
		done = true
	}
	fieldLine := fieldLines[0]
	parsedBytes := len(fieldLine) 
	if len(fieldLine) > 0 {
		if fieldLine[0] == ' ' {
			return 0, done, errors.New("whitespace not allowed at beginning of line")
		}
		fieldKey, fieldValue, _ := strings.Cut(fieldLine, ":")
		if fieldKey[len(fieldKey)-1] == ' ' {
			return 0, done, errors.New("whitespace not allowed between key and ':'")
		}
		err = checkFieldName(fieldKey)
		if err != nil {
			return 0, done, err
		}
		fieldValue = strings.Trim(fieldValue, " ")
		fieldKey = strings.ToLower(fieldKey)
		old, exists := h[fieldKey]
		if exists {
			h[fieldKey] = old + "; " + fieldValue
		} else {
			h[fieldKey] = fieldValue
		}
		fmt.Printf("Parsed sequence %s (%d)\n", fieldLine, parsedBytes)
	}
	return parsedBytes + 2, done, nil
}
