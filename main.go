package main

import (
	"fmt"
	"os"
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
	
	var BytesRead = ChunkLengh
	chunk := make([]byte, ChunkLengh)
	for BytesRead == ChunkLengh { 
		BytesRead, err := f.Read(chunk)
		check(err)
		fmt.Printf("read: %s\n", string(chunk[:BytesRead]))
	}
	// var data, err =  os.Open("./message.txt")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(data)
}
