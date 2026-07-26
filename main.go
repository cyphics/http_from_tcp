package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		if errors.Is(e, io.EOF) {
			// fmt.Print("EOF")
		} else {
			fmt.Print(e)
			os.Exit(1)
		}
	}
}

func read_file(f io.ReadCloser, str_chan chan<- string) {
	defer f.Close()
	defer close(str_chan)
	const ChunkLenght = 8
	var current_line = ""
	
	chunk := make([]byte, ChunkLenght)
	for { 
		BytesRead, err := f.Read(chunk)
		check(err)
		current_line += string(chunk[:BytesRead])
		var split = strings.Split(current_line, "\n")
		if len(split) > 1 {
			str_chan<-split[0]
			// fmt.Printf("load %s\n", split[0])
			current_line = split[1]
		}

		if BytesRead < ChunkLenght {
			break
		}
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	var chan_str = make(chan string, 1)
	go read_file(f, chan_str)
	return chan_str
}

func main() {
	f, err := os.Open("message.txt")
	check(err)
	var lines = getLinesChannel(f)
	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}
}
