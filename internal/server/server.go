package server

import (
	"bytes"
	"fmt"
	"http_from_tcp/internal/request"
	"http_from_tcp/internal/response"
	"io"
	"log"
	"net"
)

type Server struct{
	listener net.Listener
	handler Handler
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}
		fmt.Println("Accepted connection from", conn.RemoteAddr())
		s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	request, err := request.RequestFromReader(conn)
	if err != nil { fmt.Printf("Error parsing request: %s\n", err) }

	writeBuffer := new(bytes.Buffer)
	handErr := s.handler(writeBuffer, request)
	err = response.WriteStatusLine(conn, handErr.StatusCode)
	if err != nil { fmt.Printf("Error writing status line: %s\n", err) }
	err = response.WriteHeaders(conn, response.GetDefaultHeaders(0))
	if err != nil { fmt.Printf("Error handling request: %s\n", err) }
	conn.Write(writeBuffer.Bytes())
	conn.Close()
}


func Serve(handler Handler, port int) (*Server, error){
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	server := Server{listener: listener, handler: handler}
	go func() {
		server.listen()
	}()
	return &server, nil
}

func (s *Server)Close() {
	log.Println("closing server")
	err := s.listener.Close()
	if err != nil {
		log.Fatalf("Error closing connection: %s\n", err)
	}
}
