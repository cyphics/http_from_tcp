package main

import (
	"crypto/sha256"
	"fmt"
	"http_from_tcp/internal/request"
	"http_from_tcp/internal/headers"
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

func basicHandler (w *response.Writer, req *request.Request) {
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

func handleChunkedBody(w *response.Writer, body *io.ReadCloser) {
	buffer := make([]byte, 1024)
	for {
		n, err := body.Read(buffer)
		if err != nil && err != io.EOF { 
			log.Fatalf("Error reading from httpbin response: %s\n", err.Error()) 
		}
		if err == io.EOF {
			w.WriteChunkedBodyDone() 
			break
		} else {
			w.WriteChunkedBody(buffer[:n])
		}
	}
}
func chunkedHandler(w *response.Writer, res *http.Response) {
	fmt.Println("Handling chunks")
	headers := response.GetDefaultHeaders(0)
	headers.Replace("content-type", "application/json")
	headers["transfer-encoding"] = "chunked"
	w.WriteHeaders(headers)
	body := res.Body
	handleChunkedBody(w, body)
}

func trailerHandler(w *response.Writer, res *http.Response) {
	fmt.Println("Handling trailers")
	heads := response.GetDefaultHeaders(0)
	heads.Replace("content-type", "text/plain")
	headers["transfer-encoding"] = "chunked"
	heads["trailer"] = "X-Content-SHA256, X-Content-Length"
	w.WriteHeaders(heads)
	buffer := make([]byte, 4096)
	r, err := res.Body.Read(buffer)
	if err != nil {
		fmt.Errorf("Error reading body: ", err)
	}
	// fmt.Printf("Writing buffer: \n%s\n", buffer)
	w.WriteBody(buffer[:r])
	sha := sha256.Sum256(buffer)
	trailers := headers.NewHeaders()
	trailers["X-Content-Length"] = fmt.Sprintf("%d", len(buffer))
	trailers["X-Content-SHA256"] = fmt.Sprintf("%x", sha)
	w.WriteTrailers(trailers)
}

func httpbingoHandler(w *response.Writer, req *request.Request) {
	target := req.RequestLine.RequestTarget
	fmt.Printf("Target: %s\n", target)
	path := strings.TrimPrefix(target, "/httpbingo")
	fmt.Printf("Path: %s\n", path)
	url := fmt.Sprintf("https://httpbingo.org/%s", path)
	fmt.Printf("URL: %s\n", url)
	res, err := http.Get(url)
	w.WriteStatusLine(response.StatusCode(res.StatusCode))
	if err != nil { log.Fatalf("Error forwarding to httpbin: %s\n", err.Error()) }
	if strings.HasPrefix(path, "/stream") {
		chunkedHandler(w, res)
	} else if strings.HasPrefix(path, "/range") {
		trailerHandler(w, res)
	} else {
		fmt.Errorf("Error: unhandled path: %s\n", path)
	}
}

func handle(w *response.Writer, req *request.Request) {
	target := req.RequestLine.RequestTarget
	if strings.HasPrefix(target, "/httpbingo") {
		httpbingoHandler(w, req)
	} else {
		basicHandler(w, req)
	}
}

func main() {
	server, err := server.Serve(handle, port)
	if err != nil { log.Fatalf("Error starting server: %v", err) }
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
