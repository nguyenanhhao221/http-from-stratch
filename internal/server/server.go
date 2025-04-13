package server

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"httpfromtcp.haonguyen.tech/internal/request"
	"httpfromtcp.haonguyen.tech/internal/response"
)

type Server struct {
	listener net.Listener
	isClosed atomic.Bool
	handler  Handler
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	s := &Server{
		listener: listener,
		handler:  handler,
	}

	go s.listen()

	return s, nil
}

func (s *Server) Close() error {
	s.isClosed.Store(true)
	return s.listener.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.isClosed.Load() {
				return
			}
			log.Printf("error when accepting connection %v", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("error closing connection in handle: %v\n", err)
		}
	}()

	r, err := request.RequestFromReader(conn)
	if err != nil {
		log.Printf("error: parsing request: %v\n", err)
		hErr := &HandlerError{
			StatusCode: response.StatusBadRequest,
			Message:    err.Error(),
		}
		hErr.Write(conn)
		return
	}
	buf := bytes.NewBuffer([]byte{})
	handlerErr := s.handler(buf, r)
	if handlerErr != nil {
		handlerErr.Write(conn)
		return
	}

	b := buf.Bytes()

	err = response.WriteStatusLine(conn, response.StatusOK)
	if err != nil {
		log.Printf("error handle connection when writing status line: %v\n", err)
		return
	}
	headers := response.GetDefaultHeaders(len(b))
	err = response.WriteHeaders(conn, headers)
	if err != nil {
		log.Printf("error handle connection when writing headers: %v\n", err)
		return
	}
	_, err = conn.Write(b)
	if err != nil {
		log.Printf("error writing body: %v\n", err)
		return
	}
}
