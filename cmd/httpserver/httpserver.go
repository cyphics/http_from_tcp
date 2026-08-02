package main

import (
	"http_from_tcp/internal/request"
	"http_from_tcp/internal/response"
	"http_from_tcp/internal/server"
	"log"
	"os"
	"os/signal"
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
	w.WriteHeaders(response.GetDefaultHeaders(len(msg)))
	w.WriteBody([]byte(msg))
}

func main() {
	server, err := server.Serve(handler, port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
