package main

import (
	"fmt"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		os.Exit(0)
	}
}

func main() {
	f, err := os.Open("message.txt")
	check(err)
	const ChunkLengh = 8

	var current_line = ""
	
	var BytesRead = ChunkLengh
	chunk := make([]byte, ChunkLengh)
	for BytesRead == ChunkLengh { 
		BytesRead, err := f.Read(chunk)
		check(err)
		current_line += string(chunk[:BytesRead])
		var split = strings.Split(current_line, "\n")
		if len(split) > 1 {
			fmt.Printf("read: %s\n", split[0])
			current_line = split[1]
		}
		// fmt.Println(len(split))
	}
}
