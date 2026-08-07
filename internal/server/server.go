package server

import (
	"fmt"
	"http_from_tcp/internal/request"
	"http_from_tcp/internal/response"
	"log"
	"net"
)

type Server struct {
	listener net.Listener
	handler  Handler
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w *response.Writer, req *request.Request)

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil { log.Fatalf("error: %s\n", err.Error()) }
		fmt.Println("Accepted connection from", conn.RemoteAddr())
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	request, err := request.RequestFromReader(conn)
	if err != nil { return }
	writeBuffer := response.Writer{}
	s.handler(&writeBuffer, request)
	conn.Write(writeBuffer.Buffer.Bytes())
}

func Serve(handler Handler, port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil { return nil, err }
	server := Server{listener: listener, handler: handler}
	go server.listen()
	return &server, nil
}

func (s *Server) Close() {
	log.Println("closing server")
	err := s.listener.Close()
	if err != nil { log.Fatalf("Error closing connection: %s\n", err) }
}
