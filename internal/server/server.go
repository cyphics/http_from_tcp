package server

import (
	"fmt"
	"http_from_tcp/internal/response"
	"log"
	"net"
)

type Server struct{
	listener net.Listener
}

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
	var err error
	err = response.WriteStatusLine(conn, response.StatusOK)
	if err != nil {
		fmt.Printf("Error writing status line: %s\n", err)
	}
	err = response.WriteHeaders(conn, response.GetDefaultHeaders(0))
	if err != nil {
		fmt.Printf("Error handling request: %s\n", err)
	}
	err = conn.Close()
	if err != nil {
		fmt.Printf("Error closing connection: %s\n", err)
	}
}


func Serve(port int) (*Server, error){
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	server := Server{listener:listener}
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
