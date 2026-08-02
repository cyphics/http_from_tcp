package main

import (
	"http_from_tcp/internal/request"
	"http_from_tcp/internal/response"
	"http_from_tcp/internal/server"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func handler (w io.Writer, req *request.Request) *server.HandlerError {
	if req.RequestLine.RequestTarget == "/yourproblem" {
		return &server.HandlerError{
			StatusCode: response.StatusBadRequest, 
			Message: "Your problem is not my problem\n",
		}
	}
	if req.RequestLine.RequestTarget == "/myproblem" {
		return &server.HandlerError{
			StatusCode: response.StatusInternalError, 
			Message: "Woopsie, my bad\n",
		}
	}

	w.Write([]byte("All good, frfr\n"))
	return &server.HandlerError{
		StatusCode: response.StatusOK, 
		Message: "",
	}
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
