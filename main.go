package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		if errors.Is(e, io.EOF) {
		} else {
			fmt.Printf("ERROR: %s", e)
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
	var listener, err = net.Listen("tcp", "127.0.0.1:42069")
	check(err)
	defer listener.Close()

	for {
		var connection, err = listener.Accept()
		fmt.Println("Connecton accepted.")
		check(err)
		var lines = getLinesChannel(connection)
		for line := range lines {
			fmt.Println(line)
		}
		fmt.Println("Connection closed.")
	}
}
