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

func (h Headers) Replace(key string, value string) {
	h[strings.ToLower(key)] = value
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

// Parse takes a chunk of bytes and attempts to parse next header line if exists
func (h Headers) Parse(data []byte) (n int, done bool, err error) {

	if strings.HasPrefix(string(data), "\r\n") {
		return 2, true, nil
	}
	fieldLine, _, found := strings.Cut(string(data), "\r\n")
	if !found {
		return 0, false, nil
	}
	// fmt.Printf("Parsing header segment \"%s\"\n", fieldLine)
	parsedBytes := len(fieldLine) 
	if len(fieldLine) > 0 {
		if fieldLine[0] == ' ' {
			return 0, false, errors.New("whitespace not allowed at beginning of line")
		}
		fieldKey, fieldValue, _ := strings.Cut(fieldLine, ":")
		if fieldKey[len(fieldKey)-1] == ' ' {
			return 0, false, errors.New("whitespace not allowed between key and ':'")
		}
		err = checkFieldName(fieldKey)
		if err != nil {
			return 0, false, err
		}
		fieldValue = strings.Trim(fieldValue, " ")
		fieldKey = strings.ToLower(fieldKey)
		old, exists := h.Get(fieldKey)
		if exists {
			h[fieldKey] = old + "; " + fieldValue
		} else {
			h[fieldKey] = fieldValue
		}
		// fmt.Printf("Parsed sequence %s (%d)\n", fieldLine, parsedBytes)
	}
	return parsedBytes+2, false, nil
}

// Get returns the value of the header key, ensuring lower case
func (h Headers) Get(key string) (string, bool)  {
	value, exists := h[strings.ToLower(key)]
	return value, exists
}
