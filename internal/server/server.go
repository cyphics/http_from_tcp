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

// func (h *HandlerError) Write(conn net.Conn) {
// 	defer conn.Close()
// 	response.WriteStatusLine(conn, h.StatusCode)
// 	response.WriteHeaders(conn, response.GetDefaultHeaders(0))
// }

type Handler func(w *response.Writer, req *request.Request)

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}
		fmt.Println("Accepted connection from", conn.RemoteAddr())
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	request, err := request.RequestFromReader(conn)
	if err != nil {
		// hErr := &HandlerError{
		// 	StatusCode: response.StatusBadRequest,
		// 	Message:    err.Error(),
		// }
		// hErr.Write(conn)
		return
	}

	writeBuffer := response.Writer{}
	s.handler(&writeBuffer, request)
	// handErr := s.handler(writeBuffer, request)
	// if handErr != nil {
	// 	handErr.Write(conn)
	// 	return
	// }
	// err = response.WriteStatusLine(conn, response.StatusOK)
	// if err != nil { fmt.Printf("Error writing status line: %s\n", err) }
	// buff := writeBuffer.Bytes()
	// err = response.WriteHeaders(conn, response.GetDefaultHeaders(len(buff)))
	// if err != nil { fmt.Printf("Error handling request: %s\n", err) }
	conn.Write(writeBuffer.Buffer.Bytes())
}

func Serve(handler Handler, port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	server := Server{listener: listener, handler: handler}
	go server.listen()
	return &server, nil
}

func (s *Server) Close() {
	log.Println("closing server")
	err := s.listener.Close()
	if err != nil {
		log.Fatalf("Error closing connection: %s\n", err)
	}
}
