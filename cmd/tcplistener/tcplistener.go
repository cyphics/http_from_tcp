package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"http_from_tcp/internal/request"
)

const port = ":42069"

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("error listening for TCP traffic: %s\n", err.Error())
	}
	defer listener.Close()

	fmt.Println("Listening for TCP traffic on", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}
		fmt.Println("Accepted connection from", conn.RemoteAddr())

		req, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Printf("Error reading request: %s\n", err.Error())
		}
		reqLine := req.RequestLine
		fmt.Printf("Requestion line:\n- Method: %s\n- Target: %s\n- Verstion: %s\n", reqLine.Method, reqLine.RequestTarget, reqLine.HTTPVersion)
		fmt.Println("Headers:")
		headers := req.Headers
		for k, v := range headers {
			fmt.Printf("- %s: %s\n", k, v)
		}
		fmt.Println("Body:")
		fmt.Println(string(req.Body))
		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)
	go func() {
		defer f.Close()
		defer close(lines)
		currentLineContents := ""
		for {
			b := make([]byte, 8, 8)
			n, err := f.Read(b)
			if err != nil {
				if currentLineContents != "" {
					lines <- currentLineContents
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				return
			}
			str := string(b[:n])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				lines <- fmt.Sprintf("%s%s", currentLineContents, parts[i])
				currentLineContents = ""
			}
			currentLineContents += parts[len(parts)-1]
		}
	}()
	return lines
}

