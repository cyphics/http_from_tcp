package main

import (
	"fmt"
	"http_from_tcp/internal/request"
	"http_from_tcp/internal/response"
	"http_from_tcp/internal/server"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const port = 42069

func handler (w *response.Writer, req *request.Request) {
	var msg string
	var status response.StatusCode
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		status = response.StatusBadRequest
		msg = `<html>
			<head>
				<title>400 Bad Request</title>
			</head>
			<body>
				<h1>Bad Request</h1>
				<p>Your request honestly kinda sucked.</p>
			</body>
		</html>`

	case "/myproblem":
		status = response.StatusInternalError
		msg = `<html>
			<head>
				<title>500 Internal Server Error</title>
			</head>
			<body>
				<h1>Internal Server Error</h1>
				<p>Okay, you know what? This one is on me.</p>
			</body>
		</html>`
	default  :
		status = response.StatusOK
		msg = `<html>
	<head>
		<title>200 OK</title>
	</head>
	<body>
		<h1>Success!</h1>
		<p>Your request was an absolute banger.</p>
	</body>
</html>`
	}
	w.WriteStatusLine(status)
	headers := response.GetDefaultHeaders(len(msg))
	headers.Replace("content-type", "text/html")
	w.WriteHeaders(headers)
	w.WriteBody([]byte(msg))
}

func chunkHanlder(w *response.Writer, req *request.Request) {
	fmt.Println("Using chunkHandler")
	target := req.RequestLine.RequestTarget
	if strings.HasPrefix(target, "/httpbin") {
		path := strings.TrimPrefix(target, "/httpbin")
		url := fmt.Sprintf("https://httpbin.org/%s", path)
		res, err := http.Get(url)
		if err != nil { log.Fatalf("Error forwarding to httpbin: %s\n", err.Error()) }
		buffer := make([]byte, 1024)
		for {
			n, err := res.Body.Read(buffer)
			if err != nil && err != io.EOF { 
				log.Fatalf("Error reading from httpbin response: %s\n", err.Error()) 
			}
			w.WriteStatusLine(http.StatusOK)
			headers := response.GetDefaultHeaders(0)
			headers.Replace("content-type", "application/json")
			headers["transfer-encoding"] = "chunked"
			w.WriteHeaders(headers)
			w.WriteChunkedBody(buffer[:n])
			fmt.Printf("buffer: \n%s\n", buffer[:n])
			if err == io.EOF {
				fmt.Printf("Last buffer\n%s\n", buffer[:n])
				fmt.Println("Write Chunked Body Done!")
				w.WriteChunkedBodyDone() 
				break
			}
		}
	} else {
		handler(w, req)
	}
	fmt.Println("Finished handling request")
}
func main() {
	server, err := server.Serve(chunkHanlder, port)
	if err != nil { log.Fatalf("Error starting server: %v", err) }
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
